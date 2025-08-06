# ✅ Checklist Tự review
## 🧹 Code sạch & format chuẩn
- [ ] Code sạch, dễ đọc, theo convention của team (naming, cấu trúc)
- [ ] Đã chạy `go fmt`, `goimports`, `golangci-lint` (nếu có)
- [ ] Không còn đoạn code tạm, `fmt.Println`, `log.Println`, comment debug
- [ ] Không còn biến, hàm, import không sử dụng
- [ ] Hàm không quá dài, có thể tách nhỏ nếu vượt quá 50 dòng
- [ ] Code không lặp lại, đã tái sử dụng hợp lý

## 🧠 Logic & Xử lý lỗi
- [ ] Logic đúng yêu cầu, không dư thừa, rõ ràng
- [ ] Xử lý đầy đủ các trường hợp `nil`, `error`, `invalid input`
- [ ] Mỗi thao tác có thể lỗi (I/O, DB, parse...) đều được check `err`
- [ ] Không để `err` bị silent — luôn log hoặc trả về (tùy ngữ cảnh)
- [ ] Sử dụng `context` đúng cách (nếu dùng), tránh leak
- [ ] Nếu là `goroutine`, có quản lý lifecycle và tránh leak

## 🌐 HTTP Handler (Gin)
- [ ] Sử dụng đúng HTTP method (GET, POST, PUT, DELETE…)
- [ ] Validate request (param, query, body) đầy đủ, rõ ràng
- [ ] Sử dụng binding (`ShouldBindJSON`, `ShouldBindUri`, v.v.) đúng cách
- [ ] Xử lý trả về `status code` phù hợp (`200`, `201`, `400`, `404`, `500`…)
- [ ] Response rõ ràng, có `message`, `data`, `error` (nếu cần)
- [ ] Không leak thông tin nhạy cảm trong response lỗi

## 🗃️ Struct & Layer
- [ ] Struct có annotation JSON đầy đủ (`json:"field_name"`)
- [ ] Sử dụng DTO/input/output rõ ràng, không lẫn lộn model DB và response
- [ ] Đã tách rõ các layer: handler → service → repository (nếu theo clean code)
- [ ] Middleware viết rõ ràng, dễ đọc, không làm chậm hệ thống

## 🧪 Unit Test & Mock
- [ ] Có test cho handler/service chính (ít nhất happy case và edge case)
- [ ] Nếu dùng DB/service ngoài, có mock hoặc test isolation
- [ ] Có test validate đầu vào (request, param, query)
- [ ] Không dùng test hardcode giá trị không ổn định (vd: time.Now() mà không mock)
- [ ] Tên test rõ ràng: `TestFunction_Scenario_Expected`

## 🔐 Security & Config
- [ ] Không hardcode secret, token, config trong code
- [ ] Không để lộ thông tin nhạy cảm qua log hoặc panic
- [ ] Validate và sanitize đầu vào người dùng
- [ ] Nếu có auth, kiểm tra middleware hoạt động đúng

## 🤖 Sử dụng GitHub Copilot
- [ ] Đã kiểm tra kỹ mọi đoạn code được gợi ý từ GitHub Copilot
- [ ] Không để lại phần comment hoặc TODO do Copilot sinh ra
- [ ] Không commit đoạn code mà chưa hiểu rõ logic hoặc mục đích
- [ ] Đảm bảo mọi function/code block từ Copilot phù hợp với context dự án
- [ ] Tối ưu lại code được gợi ý nếu cần — tránh lạm dụng hoặc copy nguyên mẫu

## Related Tickets
- ticket redmine link

## WHAT (optional)
- Change number items `completed/total` in admin page.

## HOW
- I edit js file, inject not_vary_normal items in calculate function.

## WHY (optional)
- Because in previous version - number just depends on `normal` items. But in new version, we have `state` and `confirm_state` depends on both `normal` + `not_normal` items.

## Evidence (Screenshot or Video)


## Notes (Kiến thức tìm hiểu thêm)
