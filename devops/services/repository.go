package services

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"devops/models"
	"devops/utils"
)

// RepositoryService 仓库服务
type RepositoryService struct {
	DB       *gorm.DB
	basePath string
}

// NewRepositoryService 创建仓库服务实例
func NewRepositoryService(db *gorm.DB) *RepositoryService {
	basePath := filepath.Join(os.TempDir(), "devops", "repositories")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		utils.Logger.Error("创建仓库目录失败", zap.Error(err))
	}
	return &RepositoryService{
		DB:       db,
		basePath: basePath,
	}
}

// Branch 分支信息
type Branch struct {
	Name   string `json:"name"`
	IsHead bool   `json:"isHead"`
}

// Commit 提交信息
type Commit struct {
	Hash    string    `json:"hash"`
	Author  string    `json:"author"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

// SyncRepository 同步仓库
func (s *RepositoryService) SyncRepository(repo *models.Repository) error {
	// 解析仓库路径
	owner, repoName := parseGitHubURL(repo.URL)
	if owner == "" || repoName == "" {
		return fmt.Errorf("无效的 GitHub 仓库 URL: %s", repo.URL)
	}

	// 创建本地仓库目录
	repoPath := filepath.Join(s.basePath, owner, repoName)
	fmt.Println(repoPath, "这是路径")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return fmt.Errorf("创建仓库目录失败: %v", err)
	}

	// 克隆或更新仓库
	var r *git.Repository
	var err error

	if _, err := os.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
		// 克隆仓库
		r, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL: repo.URL,
			Auth: &gitHttp.BasicAuth{
				Username: "token",
				Password: repo.Token,
			},
		})
		if err != nil {
			return fmt.Errorf("克隆仓库失败: %v", err)
		}
	} else {
		// 打开现有仓库
		r, err = git.PlainOpen(repoPath)
		if err != nil {
			return fmt.Errorf("打开仓库失败: %v", err)
		}

		// 拉取更新
		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("获取工作目录失败: %v", err)
		}

		err = w.Pull(&git.PullOptions{
			Auth: &gitHttp.BasicAuth{
				Username: "token",
				Password: repo.Token,
			},
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("拉取更新失败: %v", err)
		}
	}

	// 获取默认分支
	head, err := r.Head()
	if err != nil {
		return fmt.Errorf("获取 HEAD 失败: %v", err)
	}

	// 更新仓库信息
	repo.DefaultBranch = head.Name().Short()
	repo.Status = "active"

	return nil
}

// GetBranches 获取分支列表
func (s *RepositoryService) GetBranches(ctx context.Context, repo *models.Repository) ([]Branch, error) {
	// 克隆仓库到内存
	r, err := git.CloneContext(ctx, memory.NewStorage(), nil, &git.CloneOptions{
		URL: repo.URL,
		Auth: &gitHttp.BasicAuth{
			Username: "token",
			Password: repo.Token,
		},
		ProxyOptions: transport.ProxyOptions{
			URL: "http://127.0.0.1:7890",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("clone repository failed: %v", err)
	}

	// 获取所有分支
	branches, err := r.Branches()
	if err != nil {
		return nil, fmt.Errorf("get branches failed: %v", err)
	}

	var result []Branch
	head, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("get head failed: %v", err)
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		branch := Branch{
			Name:   ref.Name().Short(),
			IsHead: ref.Hash() == head.Hash(),
		}
		result = append(result, branch)
		return nil
	})

	return result, err
}

// GetCommits 获取提交历史
func (s *RepositoryService) GetCommits(ctx context.Context, repo *models.Repository, branch string) ([]Commit, error) {
	// 克隆仓库到内存
	r, err := git.CloneContext(ctx, memory.NewStorage(), nil, &git.CloneOptions{
		URL: repo.URL,
		Auth: &gitHttp.BasicAuth{
			Username: "token",
			Password: repo.Token,
		},
		ProxyOptions: transport.ProxyOptions{
			URL: "http://127.0.0.1:7890",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("clone repository failed: %v", err)
	}

	// 获取分支引用
	var ref *plumbing.Reference
	if branch == "" {
		ref, err = r.Head()
	} else {
		ref, err = r.Reference(plumbing.NewBranchReferenceName(branch), true)
	}
	if err != nil {
		return nil, fmt.Errorf("get branch reference failed: %v", err)
	}

	// 获取提交历史
	commitIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("get commit history failed: %v", err)
	}

	var commits []Commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		commit := Commit{
			Hash:    c.Hash.String(),
			Author:  c.Author.Name,
			Message: c.Message,
			Date:    c.Author.When,
		}
		commits = append(commits, commit)
		return nil
	})

	return commits, err
}

// GetFiles 获取文件列表
func (s *RepositoryService) GetFiles(repo *models.Repository, path string) ([]string, error) {
	// 解析仓库路径
	owner, repoName := parseGitHubURL(repo.URL)
	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("无效的 GitHub 仓库 URL: %s", repo.URL)
	}

	// 打开仓库
	repoPath := filepath.Join(s.basePath, owner, repoName)
	_, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("打开仓库失败: %v", err)
	}

	// 读取目录内容
	dir := filepath.Join(repoPath, path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name())
	}

	return fileNames, nil
}

// parseGitHubURL 解析 GitHub URL
func parseGitHubURL(url string) (owner, repo string) {
	// 移除 URL 前缀和后缀
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimSuffix(url, ".git")

	// 分割获取 owner 和 repo
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
