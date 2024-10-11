class CookieService {
    setCookie(name, value, days = 7, path = '/') {
      const expires = new Date(Date.now() + days * 864e5).toUTCString();
      const stringValue = typeof value === 'object' ? JSON.stringify(value) : String(value);
      document.cookie = `${name}=${encodeURIComponent(stringValue)}; expires=${expires}; path=${path}`;
    }
  
    getCookie(name) {
      const value = document.cookie.split('; ').reduce((r, v) => {
        const parts = v.split('=');
        return parts[0] === name ? decodeURIComponent(parts[1]) : r;
      }, '');
      try {
        return JSON.parse(value);
      } catch {
        return value;
      }
    }
  
    deleteCookie(name, path = '/') {
      this.setCookie(name, '', -1, path);
    }
  }
  
  export const cookieService = new CookieService();