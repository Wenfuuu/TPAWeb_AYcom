import axios from "axios";
import { fetchUser } from "./user";

export type Conversation = {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
  is_deleted: boolean;
}

export type Message = {
  id: string;
  conversation_id: string;
  sender_id: string;
  content: string;
  media_type: "text" | "image" | "gif" | "video";
  media_url: string;
  created_at: string;
  updated_at: string;
  is_deleted: boolean;
  is_unsent: boolean;
  profile_picture: string;
};

const BASE_API_URL = "http://localhost:8080/api/chat";

let messages: Message[] = [];

export async function fetchUserConversation(user_id: string) {
    const response = await axios.get(BASE_API_URL + "/conversations", {
        params: {
            user_id
        }
    });

    const conversations = response.data.conversations;
    return conversations;
}

export async function fetchMessageWithUser(conversation_id: string) {
    const allMessages = await fetchMessage(conversation_id);
    const allUsers = await fetchUser();

    if(!allMessages) return [];

    messages = allMessages.map(message => {
        const user = allUsers.find(u => u.id === message.sender_id);
        const profile_picture = user ? user.profile_picture : "";

        return {
            ...message,
            profile_picture
        };
    });

    return messages;
}

export async function fetchMessage(conversation_id: string) {
    const response = await axios.get(
      BASE_API_URL + "/messages",
      {
        params: {
          conversation_id,
        },
      }
    );

    const messages: Message[] = response.data.messages;
    return messages;
}