import http from '../http-common'

class EmailService {
  get(key: string) {
    return http.get(`/emails/${key}`)
  }
}

export default new EmailService()