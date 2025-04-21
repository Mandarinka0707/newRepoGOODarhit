// src/services/chatService.js

const API_URL = 'http://localhost:50051'; // Замените на URL вашего chat_service (gRPC или REST)

// Функция для получения истории сообщений (если используется REST API)
export const getChatHistory = async () => {
  try {
    const response = await fetch(`${API_URL}/messages`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Could not fetch chat history:", error);
    throw error;
  }
};

// Функция для отправки сообщений (если используется REST API)
export const sendMessage = async (message) => {
  try {
    const response = await fetch(`${API_URL}/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(message),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Could not send message:", error);
    throw error;
  }
};

// Если используется gRPC (пример - требует библиотеки grpc-web):
/*
import { ChatServiceClient } from './grpc/ChatServiceClient'; // Путь к вашему сгенерированному gRPC клиенту
import { SendMessageRequest, GetChatHistoryRequest } from './grpc/chat_pb'; // Пути к сообщениям protobuf

const chatClient = new ChatServiceClient(API_URL);

export const getChatHistory = async () => {
    return new Promise((resolve, reject) => {
        chatClient.getChatHistory(new GetChatHistoryRequest(), {}, (err, response) => {
            if (err) {
                console.error("Could not fetch chat history:", err);
                reject(err);
                return;
            }
            resolve(response.getMessagesList());
        });
    });
};

export const sendMessage = async (message) => {
    const request = new SendMessageRequest();
    request.setText(message.text); // Предполагается, что сообщение имеет поле text

    return new Promise((resolve, reject) => {
        chatClient.sendMessage(request, {}, (err, response) => {
            if (err) {
                console.error("Could not send message:", err);
                reject(err);
                return;
            }
            resolve(response); // Возвращаем ответ от сервера
        });
    });
};
*/

// WebSocket - обрабатывается напрямую в ChatRoom.jsx
// (Код WebSocket уже предоставлен в предыдущем примере)

export default { getChatHistory, sendMessage };