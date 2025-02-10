# OTP Authentication Service

## 🚀 Introduction
This project is a secure OTP-based authentication system built using Go (Gin framework), MySQL, and Redis. It provides a seamless authentication process using mobile numbers and OTPs, ensuring a smooth and secure login experience.

---

## 🎯 Features
✅ OTP-based authentication using Redis for secure storage.  
✅ User registration with mobile number verification.  
✅ Secure JWT token generation for authentication.  
✅ OTP cooldown management to prevent abuse.  
✅ OTP expiration after successful login.  
✅ Exception handling for better API reliability.  
✅ Two-step authentication during login (OTP request and OTP verification).  

---

## 🏗️ Tech Stack
- **Go (Gin framework)** - Lightweight web framework.
- **MySQL** - Database for user management.
- **Redis** - OTP storage and cooldown tracking.
- **JWT (JSON Web Tokens)** - Secure authentication.

---

## 📌 API Endpoints
### 1️⃣ Register User
**Endpoint:** `POST /register`  
Registers a new user and sends an OTP for verification.

#### Request Body:
```json
{
  "mobile": "9876543210",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com"
}
```
#### Response:
```json
{
  "message": "User registered successfully. OTP sent!"
}
```

### 2️⃣ Verify OTP
**Endpoint:** `POST /verify-otp`  
Verifies the OTP sent to the user.

#### Request Body:
```json
{
  "mobile": "9876543210",
  "otp": "123456"
}
```
#### Response:
```json
{
  "message": "OTP verified successfully! Account activated."
}
```

### 3️⃣ Login User
**Endpoint:** `POST /login`  
Logs in the user using OTP authentication.

#### Request Body:
```json
{
  "mobile": "9876543210",
  "otp": "123456"
}
```
#### Response:
```json
{
  "token": "<jwt_token>",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "login_at": 1739185527
}
```

#### Login Flow:
1️⃣ **User initiates login with a mobile number** → OTP is sent as a 2FA step.  
2️⃣ **User enters the received OTP** → System verifies the OTP.  
3️⃣ **If the user is not verified** → System returns `Account is not verified`.  
4️⃣ **If OTP is valid and user is verified** → System issues JWT token for authentication.  
5️⃣ **OTP is invalidated after successful use** to prevent replay attacks.  

---

## 🔄 API Flow
1️⃣ **User registers** → OTP is sent for verification.  
2️⃣ **User verifies OTP** → Account is activated.  
3️⃣ **User logs in** → OTP is verified and JWT token is issued.  
4️⃣ **OTP is invalidated after use** → Prevents replay attacks.  

---

## 🔐 Security & Exception Handling
- OTP cooldown prevents multiple requests within **30 seconds**.
- OTP expires **in 5 minutes** if unused.
- **Invalid OTP** returns `401 Unauthorized`.
- **User not registered** returns `404 Not Found`.
- **Account not verified** prevents login.
- **JWT tokens** are signed securely for authentication.

---

## 📦 Setup Instructions
1️⃣ Clone the repo:  
   ```sh
   git clone https://github.com/your-repo/otp-auth.git
   ```
2️⃣ Install dependencies:  
   ```sh
   go mod tidy
   ```
3️⃣ Set up **MySQL & Redis** and configure `.env`:  
   ```
    PORT=8080
    DB_USER=root
    DB_PASSWORD=Snippet@1
    DB_NAME=otp_auth
    DB_HOST=localhost
    DB_PORT=3306
    REDIS_HOST=localhost
    REDIS_PORT=6379
    SECRET_KEY=django-insecure-2oj=x%p4nn(hvhqhem5n89iktidgdd0ma7mrs-i$8r3ushal60
    OTP_EXPIRY=5
    TWILIO_ACCOUNT_SID='AC26b7be6e305e49310f85bf7e65d08cb3'
    TWILIO_AUTH_TOKEN='eca9c34d20cf554f3c9e739a3fc29b35'
    TWILIO_PHONE_NUMBER='+14147101993'
   ```
4️⃣ Run the project:  
   ```sh
   go run main.go
   ```

---

## 🎯 Future Improvements
- **Rate limiting** to prevent spam.
- **Multi-factor authentication (MFA)** support.
- **Better logging and monitoring.**

---


