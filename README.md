# Getting Started

Run command ``go run main.go`` to start application

# API List

1. `POST /api/login`: login to application  
   How to test  
   endpoint `http://localhost:8080/api/login`  
   request body example  
       ```{
       "username": "tester01",
       "password": "1111"
       } ```

   response body example  
       ```{
       "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2ODE4MjUzMzMsImp0aSI6IjEiLCJpYXQiOjE2ODE3Mzg5MzN9.QFvBHCLFQ-NJBkChBB6PvjDXPCK-GXdbbqKzyDFQDho"
       } ```  


2. `POST /api/initialize`: initialize table  
   How to test  
   endpoint `http://localhost:8080/api/initialize`  
   --header Authorization: Bearer <access_token>  
   request body example  
       ```{
            "numTables": 2
       }```  
   response body example   
       ```{
       "success": true,
        "message": "Tables initialized successfully"
       } ```  

3. `POST /api/reserve`: reserve table
   How to test  
   endpoint `http://localhost:8080/api/reserve`  
   --header Authorization: Bearer <access_token>  
   request body example  
       ```{
       "numCustomers": 8
       } ```  
   response body example   
       ```{
       "success": true,
       "data": {
        "bookingID": 1,
        "numBookedTables": 2,
        "numRemaining": 0
         }
      } ```  

4. `POST /api/cancel`: cancel reserve table
   How to test  
   endpoint example `http://localhost:8080/api/cancel`  
   --header Authorization: Bearer <access_token>  
   request body example  
       ```{
         "bookingID": 1
       }```  
   response body example   
       ```{
   "success": true,
   "data": {
   "numFreedTables": 1,
   "numRemaining": 3
   }
   }
    } ```  

# High-Level Design
The application is designed using the Hexagonal Architecture pattern, which provides a clean separation between the business logic and the infrastructure code.

Overall, this design provides a clean separation between the different components of the application, making it easy to test, maintain, and extend.


# Stack
1. Gin as the web framework