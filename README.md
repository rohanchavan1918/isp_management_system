# isp_management_system
A ISP Management system to manage the users of ISP,and the plans.
This is a simple practice project to polish golang skills, and have a templated code which can be reused in future.


#### Completed Features
1. Login
2. Registeration
3. Authentication / Authorization with JWT
4. User Reset password.
5. Forgot password with OTP to mail.
6. Admin User can create/update/delete plans which regular users can view.
    As these plans will be common for all users they are cached (Redis) , when updated or deleted cache is also updated.


#### TODO Features
1. WRITE TEST CASES.
2. Soft delete (save is_active as false ) .
3. Other endpoints for CRUD of plans, assign plan to users, CRUD for user complaints, etc. 
4. API documentation with SWAGGER.