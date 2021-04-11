APP_NAME=zalopay-api
echo "Starting..."
echo "Docker building..."
docker build -t ${APP_NAME} -f ./Dockerfile .
echo "Docker saving..."
docker save -o ${APP_NAME}.tar ${APP_NAME}

echo "Done"
