# Mini-Twitter

Mini-Twitter is a social media platform built in Go using the Gin framework. It provides functionality for users to create accounts, post tweets, follow other users, and receive real-time notifications.

## Features

1. **Email Verification**: Utilizes Redis for storing verification codes sent to user emails during registration.
2. **PostgreSQL Database**: Main database for storing user data, tweets, and relationships. Includes indexing for faster data retrieval.
   ```sql
   CREATE INDEX idx_users_username ON users (username);
   CREATE INDEX idx_users_name ON users (name);
3. **AWS S3 Integration**: Used for saving user-uploaded photos and files.
4. **Real-time Notifications**: Implemented using Apache Kafka with WebSocket support.
5. **Load Testing**: Conducted using k6, with load test scripts for GET and POST requests.
6. **API Documentation**: Swagger documentation for easy reference and testing.
7. **Role-Based Access Control**: Utilizes Casbin for role checking with JWT. Roles include ```unauthorized```, ```user```, and ```admin```.
8. **Rate Limiting**: Middleware implemented to limit request rates.
9. **Docker Support**: Dockerfile and Docker Compose configurations for containerized deployment.

# Getting Started
## Prerequisites
 * Docker and Docker Compose
 * Go 1.20 or higher

## Environment Variables
Create a ```.env``` file in the root directory of the project using the following template
  ```bash
  # Application configuration
  APP=mini-twitter
  ENVIRONMENT=develop
  LOG_LEVEL=debug
  CONTEXT_TIMEOUT=7s
  GIN_MODE=debug

  # Server configuration
  SERVER_HOST=twitter_app
  SERVER_PORT=:7777
  SERVER_READ_TIMEOUT=10s
  SERVER_WRITE_TIMEOUT=10s
  SERVER_IDLE_TIMEOUT=120s

  # PostgreSQL database configuration
  POSTGRES_HOST=twitter_postgres
  POSTGRES_PORT=5432
  POSTGRES_USER=postgres
  POSTGRES_PASSWORD=your_postgres_password
  POSTGRES_DATABASE=twitter_db
  POSTGRES_SSL_MODE=disable

  # Redis configuration
  REDIS_HOST=twitter_redis
  REDIS_PORT=6379

  # Kafka configuration
  KAFKA_BROKER=broker:29092
  KAFKA_TOPIC=notification

  # JWT configuration
  SIGNING_KEY=your_signing_key
  ACCESS_TTL=6h
  REFRESH_TTL=24h

  # AWS S3 configuration
  AWS_ACCESS_KEY_ID=your_aws_access_key_id
  AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
  AWS_BUCKET_NAME=your_s3_bucket_name
  AWS_REGION=your_aws_region

  # Casbin authorization configuration
  CSV_FILE_PATH=./config/auth.csv
  CONF_FILE_PATH=./config/auth.conf

  # SMTP (Email) configuration
  SMTP_HOST=smtp.gmail.com
  SMTP_PORT=587
  SMTP_EMAIL=your_smtp_email
  SMTP_PASSWORD=your_smtp_password
  ```


## Running the Project
1. **Clone the repository**:
  ```bash
  git clone https://github.com/dostonshernazarov/mini-twitter.git
  cd mini-twitter
  ```

2. **Build and start the application using Docker Compose**:
  ```bash
  docker-compose up --build
  ```

3. **Access the application**:
  * The API is accessible at http://localhost:7777.
  * Swagger documentation is available at http://localhost:7777/v1/swagger/index.html.

## Load Testing with k6
Load tests can be run using ```k6```. Ensure that the ```k6``` service is istalled in device.

# Contributing
Contributions are welcome! Please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
