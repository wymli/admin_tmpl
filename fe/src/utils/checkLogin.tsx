export default function checkLogin() {
  return true // 无需登录
  return localStorage.getItem('userStatus') === 'login';
}
