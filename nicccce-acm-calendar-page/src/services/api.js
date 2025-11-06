import axios from 'axios';

// 创建axios实例
const apiClient = axios.create({
  baseURL: 'https://分形黄昏.nicccce.xyz/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    console.log(`API请求: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    // 根据后端响应格式调整
    if (response.data && response.data.code === 0) {
      return response.data.data || response.data;
    }
    return response.data;
  },
  (error) => {
    console.error('API请求错误:', error);
    return Promise.reject(error);
  }
);

// API服务类
class ApiService {
  // 获取比赛列表
  async getContests(params = {}) {
    try {
      const response = await apiClient.get('/contests', { params });
      // 返回实际的数据而不是整个响应对象
      return response.data;
    } catch (error) {
      console.error('获取比赛列表失败:', error);
      throw error;
    }
  }

  // 根据ID获取比赛
  async getContestById(id) {
    try {
      const response = await apiClient.get(`/contests/${id}`);
      return response.data;
    } catch (error) {
      console.error('获取比赛详情失败:', error);
      throw error;
    }
  }

  // 根据平台获取比赛
  async getContestsByPlatform(platform) {
    try {
      const response = await apiClient.get(`/contests/platform/${platform}`);
      return response.data;
    } catch (error) {
      console.error('根据平台获取比赛失败:', error);
      throw error;
    }
  }

  // 根据状态获取比赛
  async getContestsByStatus(status) {
    try {
      const response = await apiClient.get(`/contests/status/${status}`);
      return response.data;
    } catch (error) {
      console.error('根据状态获取比赛失败:', error);
      throw error;
    }
  }

  // 刷新所有平台数据
  async refreshAllPlatforms() {
    try {
      const response = await apiClient.post('/refresh');
      return response.data;
    } catch (error) {
      console.error('刷新所有平台失败:', error);
      throw error;
    }
  }

  // 刷新单个平台数据
  async refreshSinglePlatform(platform) {
    try {
      const response = await apiClient.post(`/refresh/${platform}`);
      return response.data;
    } catch (error) {
      console.error('刷新单个平台失败:', error);
      throw error;
    }
  }

  // 获取刷新状态
  async getRefreshStatus() {
    try {
      const response = await apiClient.get('/refresh/status');
      return response;
    } catch (error) {
      console.error('获取刷新状态失败:', error);
      throw error;
    }
  }

  // 获取速率限制信息
  async getRateLimitInfo(platform = 'all') {
    try {
      const response = await apiClient.get('/refresh/status');
      return response;
    } catch (error) {
      console.error('获取速率限制信息失败:', error);
      throw error;
    }
  }

  // 获取比赛统计信息
  async getContestStats() {
    try {
      const response = await apiClient.get('/refresh/status');
      return response;
    } catch (error) {
      console.error('获取刷新状态失败:', error);
      throw error;
    }
  }

  // 获取刷新日志
  async getRefreshLogs(limit = 50) {
    try {
      const response = await apiClient.get('/admin/contests/logs', { 
        params: { limit } 
      });
      return response;
    } catch (error) {
      console.error('获取刷新日志失败:', error);
      throw error;
    }
  }
}

// 创建API服务实例
const apiService = new ApiService();

export default apiService;