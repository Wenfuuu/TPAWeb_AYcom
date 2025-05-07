import axios from "axios";

export type User = {
    id: string
    name: string
    username: string
    email: string
    password: string
    gender: string
    profile_picture: string
    banner: string
    date_of_birth: string
    is_active: boolean
    is_banned: boolean
    role: string
    is_subscribed: boolean
    is_premium: boolean
    bio: string
    created_at: string
    is_followed: boolean
}

export type UserThread = {
    name: string;
    username: string;
    isPremium: boolean;
    profile: string;
};

type Follows = {
  following_id: string;
  follower_id: string;
}

const BASE_API_URL = "http://localhost:8080/api/users";
export async function fetchUser() {
    const response = await axios.get<User[]>(BASE_API_URL + "/user");
    const users = response.data;
    return users;
}

export async function followUser(follower_id: string, following_id: string) {
  console.log(follower_id);
  console.log(following_id);

  await axios.post(BASE_API_URL + "/follow", {
    follower_id,
    following_id
  });
}

export async function fetchFollowing(user_id: string) {
  const response = await axios.get<Follows[]>(BASE_API_URL + "/following", {
    params: {
      user_id
    }
  });
  const users = response.data;
  return users;
}

export async function fetchProfileUser(id: string) {
  const response = await axios.get<User>(BASE_API_URL + "/user-profile", {
    params: {
      id,
    },
  });
  const user = response.data;
  return user;
}

export async function fetchWhoToFollow(user_id: string) {
    const response = await axios.get<User[]>(BASE_API_URL + "/who-to-follow", {
        params: {
            user_id
        }
    });
    const toFollows = response.data;
    const followings = await fetchFollowing(user_id);
    const followingIds = new Set(followings.map((user) => user.following_id));
    console.log(followingIds);

    const users = toFollows.map((user) => {
      return {
        ...user,
        is_followed: followingIds.has(user.id),
      };
    });

    return users;
}

export async function registerPremium(userId: string, identityNumber: string, facePicture: string, reason: string): Promise<string | boolean> {
  try {
    await axios.post(BASE_API_URL + "/register-premium", {
      userId,
      identityNumber,
      reason,
      facePicture,
    });
  } catch (error: unknown) {
    if (axios.isAxiosError(error)) {
      return error.message;
    }

    return "An unexpected error occurred.";
  }

  return true;
}