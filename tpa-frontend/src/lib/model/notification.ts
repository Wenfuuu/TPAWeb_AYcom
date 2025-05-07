import axios from "axios";
import { fetchUser } from "./user";

export type Notif = {
  id: string;
  user_id: string;
  sender_id: string;
  thread_id: string;
  content: string;
  category: string;
  created_at: string;
  profile_picture: string;
};

const BASE_API_URL = "http://localhost:8080/api/chat";
let notifs: Notif[] = [];

export async function fetchAllNotifications(user_id: string) {
    const response = await axios.get(BASE_API_URL + "/notifications", {
        params: {
            user_id
        }
    });

    const notifications: Notif[] = response.data.notifications;
    if(!notifications) return [];

    const allUsers = await fetchUser();
    notifs = notifications.map(notif => {
        // console.log(notif.sender_id);
        const user = allUsers.find((u) => u.id === notif.sender_id);
        const profile_picture = user ? user.profile_picture : "";

        return {
            ...notif,
            profile_picture
        };
    });

    return notifs;
}

export async function fetchMentionNotifications(user_id: string) {
    const response = await axios.get(BASE_API_URL + "/mention-notifications", {
        params: {
            user_id
        }
    });

    const notifications: Notif[] = response.data.notifications;
    if (!notifications) return [];



    return notifications;
}