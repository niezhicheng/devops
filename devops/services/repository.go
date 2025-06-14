package services

import (
	"devops/models"
	"devops/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"go.uber.org/zap"
)

// RepositoryService 仓库服务
type RepositoryService struct {
	basePath string
}

// NewRepositoryService 创建仓库服务实例
func NewRepositoryService() *RepositoryService {
	basePath := filepath.Join(os.TempDir(), "devops", "repositories")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		utils.Logger.Error("创建仓库目录失败", zap.Error(err))
	}
	return &RepositoryService{
		basePath: basePath,
	}
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
			Auth: &http.BasicAuth{
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
			Auth: &http.BasicAuth{
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
	//repo.LastSyncAt = time.Now()
	repo.Status = "active"

	return nil
}

// GetBranches 获取分支列表
func (s *RepositoryService) GetBranches(repo *models.Repository) ([]string, error) {
	// 解析仓库路径
	owner, repoName := parseGitHubURL(repo.URL)
	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("无效的 GitHub 仓库 URL: %s", repo.URL)
	}

	// 打开仓库
	repoPath := filepath.Join(s.basePath, owner, repoName)
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("打开仓库失败: %v", err)
	}

	// 获取所有分支
	branches, err := r.Branches()
	if err != nil {
		return nil, fmt.Errorf("获取分支列表失败: %v", err)
	}

	var branchNames []string
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		branchNames = append(branchNames, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("遍历分支失败: %v", err)
	}

	return branchNames, nil
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
