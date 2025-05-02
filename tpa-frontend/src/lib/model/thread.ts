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
  isLiked: boolean;
  isReposted: boolean;
  isBookmarked: boolean;
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

export async function fetchThreadWithUser(page: number, userid: string): Promise<(Thread & { userData: UserThread })[] | null> {
    const allThreads = await fetchThread(page);
    const allUsers = await fetchUser();
    const fileURLS = await fetchFiles();
    const likedThreads = (await fetchLikedThread(userid)) || [];
    const bookmarkThreads = (await fetchBookmarkThread(userid)) || [];

    if(!allThreads) return null;

    const likedThreadIds = new Set(likedThreads.map(t => t.id));
    const bookmarkThreadIds = new Set(bookmarkThreads.map(t => t.id));

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
            profile: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//kagari.jpg"
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
            likes: 0,
        };
        return {
            ...thread,
            created_at: formatTweetDate(new Date(thread.created_at)),
            images,
            insights,
            userData,
            isLiked: likedThreadIds.has(thread.id),
            isReposted: false,
            isBookmarked: bookmarkThreadIds.has(thread.id)
        };
    });

    return threads;
}

export async function fetchFollowingThread(page: number, userid: string): Promise<(Thread & { userData: UserThread })[] | null> {
    const allThreads = await fetchFollowThread(page, userid);
    const allUsers = await fetchUser();
    const fileURLS = await fetchFiles();
    const likedThreads = (await fetchLikedThread(userid)) || [];
    const bookmarkThreads = (await fetchBookmarkThread(userid)) || [];

    if(!allThreads) return null;

    const likedThreadIds = new Set(likedThreads?.map(t => t.id));
    const bookmarkThreadIds = new Set(bookmarkThreads.map(t => t.id));

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
            isPremium: true,
            profile: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//kagari.jpg"
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
            likes: 0,
        };
        return {
            ...thread,
            created_at: formatTweetDate(new Date(thread.created_at)),
            images,
            insights,
            userData,
            isLiked: likedThreadIds.has(thread.id),
            isReposted: false,
            isBookmarked: bookmarkThreadIds.has(thread.id)
        };
    });

    return threads;
}

export async function fetchThreadDetailWithUser(page: number, threadid: string, userid: string): Promise<(Thread & { userData: UserThread })[] | null> {
    const allThread = await fetchThreadDetail(threadid);
    // const allReplies = await fetchThreadReplies(page, threadid);
    const allUsers = await fetchUser();
    const fileURLS = await fetchFiles();
    const likedThreads = (await fetchLikedThread(userid)) || [];
    const bookmarkThreads = (await fetchBookmarkThread(userid)) || [];

    if (!allThread) return null;

    const threadsArray = Array.isArray(allThread) ? allThread : [allThread];

    const likedThreadIds = new Set(likedThreads.map(t => t.id));
    const bookmarkThreadIds = new Set(bookmarkThreads.map(t => t.id));

    const result = threadsArray.map(thread => {
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
                  isPremium: true,
                  profile: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//kagari.jpg",
              };

        let images: string[] = [];
        if (fileURLS) {
            images = fileURLS
                .filter(file => file.thread_id === thread.id)
                .map(file => file.file);
        }

        const insights = {
            replies: 0,
            reposts: 0,
            likes: 0,
        };

        return {
            ...thread,
            created_at: formatTweetDate(new Date(thread.created_at)),
            images,
            insights,
            userData,
            isLiked: likedThreadIds.has(thread.id),
            isReposted: false,
            isBookmarked: bookmarkThreadIds.has(thread.id)
        };
    });

    return result;
}

export async function fetchLikedThread(user_id: string) {
    const response = await axios.get<Thread[]>(BASE_API_URL + "/like", {
        params: {
            user_id
        }
    });
    const threads = response.data;
    return threads;
}

export async function fetchBookmarkThread(user_id: string) {
    const response = await axios.get<Thread[]>(BASE_API_URL + "/bookmark", {
        params: {
            user_id
        }
    });
    const threads = response.data;
    return threads;
}

export async function fetchFollowThread(page: number, user_id: string) {
    const response = await axios.get<Thread[]>(BASE_API_URL + "/thread-following", {
        params: {
            page,
            user_id
        }
    });
    const threads = response.data;
    return threads;
}

export async function fetchThreadReplies(page: number, thread_id: string) {
    const response = await axios.get<Thread[]>(BASE_API_URL + "/thread-reply", {
        params: {
            page,
            thread_id
        }
    });
    const threads = response.data;
    return threads;
}

export async function fetchThreadDetail(thread_id: string) {
    const response = await axios.get<Thread>(BASE_API_URL + "/thread-detail", {
        params: {
            thread_id
        }
    });
    const threads = response.data;
    return threads;
}

export async function fetchThread(page: number) {
    const response = await axios.get<Thread[]>(BASE_API_URL + "/thread", {
        params: {
            page
        }
    });
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

export async function likeThread(userId: string, threadId: string) {
    await axios.post(BASE_API_URL + "/like", {
        userId,
        threadId
    });
}

export async function unlikeThread(userId: string, threadId: string) {
    await axios.post(BASE_API_URL + "/unlike", {
        userId,
        threadId
    });
}

export async function bookmarkThread(userId: string, threadId: string) {
    await axios.post(BASE_API_URL + "/bookmark", {
        userId,
        threadId
    });
}

export async function unbookmarkThread(userId: string, threadId: string) {
    await axios.post(BASE_API_URL + "/unbookmark", {
        userId,
        threadId
    });
}