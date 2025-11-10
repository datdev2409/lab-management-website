package handlers

const (
	AUTH_INVALID_PAYLOAD_ERROR     = "Tên đăng nhập hoặc mật khẩu không hợp lệ"
	AUTH_USER_EXISTS_ERROR         = "Người dùng đã tồn tại"
	AUTH_USER_NOT_FOUND_ERROR      = "Người dùng không tồn tại"
	AUTH_LOGIN_FAILED_ERROR        = "Đăng nhập thất bại. Vui lòng kiểm tra lại tên đăng nhập và mật khẩu"
	AUTH_INVALID_ADMIN_TOKEN_ERROR = "Mã quản trị không hợp lệ. Vui lòng liên hệ quản trị viên"

	INVALID_REQUEST_PAYLOAD_ERROR = "Dữ liệu yêu cầu không hợp lệ"

	DUPLICATE_PATIENT_ERROR = "Bệnh nhân với tên và số điện thoại đã tồn tại"
	DUPLICATE_DOCTOR_ERROR  = "Bác sĩ với tên và số điện thoại đã tồn tại"
	DOCTOR_NOT_FOUND_ERROR  = "Không tìm thấy bác sĩ chỉ định"

	PATIENT_ALREADY_EXISTS = "Bệnh nhân với tên và số điện thoại đã tồn tại"
)
