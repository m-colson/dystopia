FROM python:3.10-alpine

WORKDIR /usr/src/app
RUN pip3 install flask requests

COPY . .

CMD ["python3", "Scheduler.py"]
