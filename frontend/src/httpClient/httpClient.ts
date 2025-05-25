import axios from 'axios'

// Create axios instance
const httpClient = axios.create({
  //Олег, можно трогать
  baseURL: 'https://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// интерсептор для авторизации каждого запроса
httpClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

httpClient.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // Здесь будет редирект на логин
          break
        case 403:
          // тут просто форбидден
          break
        case 404:
          // тут редирект на нот фаунд
          break
        case 500:
          // тут редирект на что-то типа упс, у нас все сломалось к хуям нахуй
          break
        default:
        // Handle other status codes
      }
    } else if (error.request) {
      console.error('No response received:', error.request)
    } else {
      console.error('Request setup error:', error.message)
    }

    return Promise.reject(error)
  },
)

export default httpClient
