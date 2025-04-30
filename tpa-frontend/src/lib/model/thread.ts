import axios from "axios";
import { fetchUser, type UserThread } from "./user";
import { format, isThisYear, differenceInHours, differenceInMinutes } from "date-fns";

export type ThreadCategory = {
    id: string
    name: string
}

export type Community = {
    id: string
    name: string
    description: string
}

export type ThreadFile = {
    id: string
    thread_id: string
    file: string
}

type ThreadInsights = {
    replies: number;
    reposts: number;
    likes: number;
};

export type Thread = {
  id: string;
  user_id: string;
  category_id: string;
  community_id: string;
  content: string;
  reply_permission: string;
  status: string;
  images: string[];
  insights: ThreadInsights;
  created_at: string;
};

const BASE_API_URL = "http://localhost:8080/api/threads";

let threads: (Thread & { userData: UserThread })[] = [];

function formatTweetDate(date: Date): string {
  const now = new Date();
  const diffInMinutes = differenceInMinutes(now, date);
  const diffInHours = differenceInHours(now, date);

  if (diffInMinutes < 60) {
    return `${diffInMinutes}m`;
  } else if (diffInHours < 24) {
    return `${diffInHours}h`;
  } else if (isThisYear(date)) {
    return format(date, "MMM d"); // e.g., "Apr 29"
  } else {
    return format(date, "MMM d, yyyy"); // e.g., "Jul 16, 2020"
  }
}

export async function fetchThreadWithUser() {
    const allThreads = await fetchThread();
    const allUsers = await fetchUser();
    const fileURLS = await fetchFiles();

    threads = allThreads.map(thread => {
        const user = allUsers.find(u => u.id === thread.user_id);
        const userData: UserThread = user
        ? {
            name: user.name,
            username: user.username,
            isPremium: user.is_premium ? true : false,
            profile: user.profile_picture,
          }
        : {
            name: "Unknown",
            username: "unknown",
            isPremium: false,
            profile: ""
        };

        let images: string[] = [];
        if(fileURLS) {
            images = fileURLS
            .filter(file => file.thread_id === thread.id)
            .map(file => file.file);
        }

        const insights: ThreadInsights = {
            replies: 0,
            reposts: 0,
            likes: 0.
        };
        return {
            ...thread,
            created_at: formatTweetDate(new Date(thread.created_at)),
            images,
            insights,
            userData,
        };
    });

    return threads;
}

export async function fetchThread() {
    const response = await axios.get<Thread[]>(BASE_API_URL + "/thread");
    const threads = response.data;
    return threads;
}

export async function fetchThreadCategory() {
    const response = await axios.get<ThreadCategory[]>(BASE_API_URL + "/thread-category");
    const categories = response.data;
    return categories;
}

export async function fetchCommunity() {
    const response = await axios.get<Community[]>(BASE_API_URL + "/community");
    const communities = response.data;
    return communities;
}

export async function fetchFiles() {
    const response = await axios.get<ThreadFile[]>(BASE_API_URL + "/thread-file");
    const files = response.data;
    return files;
}