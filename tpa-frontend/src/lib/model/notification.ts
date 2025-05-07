import axios from "axios";

export type Notification = {
  id: string;
  user_id: string;
  sender_id: string;
  thread_id: string;
  content: string;
  category: string;
  created_at: string;
};

const BASE_API_URL = "http://localhost:8080/api/chat";

export async function fetchNotifications(user_id: string) {
    const response = await axios.get(BASE_API_URL + "/notifications", {
        params: {
            user_id
        }
    });

    const notifications: Notification[] = response.data.notifications;
    return notifications;
}