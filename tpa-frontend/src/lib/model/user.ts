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
    createdAt: string
}

export type UserThread = {
    name: string;
    username: string;
    isPremium: boolean;
    profile: string;
};

const BASE_API_URL = "http://localhost:8080/api/users";
export async function fetchUser() {
    const response = await axios.get<User[]>(BASE_API_URL + "/user");
    const users = response.data;
    return users;
}

export async function fetchWhoToFolllow(user_id: string) {
    const response = await axios.get<User[]>(BASE_API_URL + "/who-to-follow", {
        params: {
            user_id
        }
    });
    const users = response.data;
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