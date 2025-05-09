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

export type SecurityQuestion = {
  user_id: string;
  security_question: string;
  security_answer: string;
}

export type PremiumRequest = {
  id: string;
  user_id: string;
  identitiy_number: string;
  reason: string;
  face_picture: string;
  created_at: string;
}

export type Reports = {
  id: string;
  user_id: string;
  reported_id: string;
  reason: string;
  created_at: string;
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
  await axios.post(BASE_API_URL + "/follow", {
    follower_id,
    following_id
  });
}

export async function unfollowUser(follower_id: string, following_id: string) {
  await axios.delete(BASE_API_URL + "/unfollow", {
    params: {
      follower_id,
      following_id
    }
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

export async function fetchUserEmailById(userId: string) {
  const user = await fetchProfileUser(userId);
  return user.email;
}

export async function setUserBanStatus(user_id: string, banned: boolean) {
  await axios.put(BASE_API_URL + "/ban", {
    user_id,
    banned,
  });
}

export async function resetPassword(user_id: string, password: string) {
  await axios.put(BASE_API_URL + "/reset-password", {
    user_id,
    password
  });
}

export async function acceptReport(report_id: string) {
  await axios.post(BASE_API_URL + "/accept-report", {
    report_id,
  });
}

export async function rejectReport(report_id: string) {
  await axios.post(BASE_API_URL + "/reject-report", {
    report_id,
  });
}

export async function fetchReports() {
  const response = await axios.get<Reports[]>(BASE_API_URL + "/reports");

  const requests = response.data;
  if(!requests) return [];
  return requests;
}

export async function acceptPremium(request_id: string) {
  await axios.post(BASE_API_URL + "/accept-premium", {
    request_id,
  });
}

export async function rejectPremium(request_id: string) {
  await axios.post(BASE_API_URL + "/reject-premium", {
    request_id,
  });
}

export async function fetchPremiumRequests() {
  const response = await axios.get<PremiumRequest[]>(BASE_API_URL + "/premium-request");

  const requests = response.data;
  if(!requests) return [];
  return requests;
}

export async function fetchSecurityQuestions(user_id: string) {
    const response = await axios.get<SecurityQuestion>(BASE_API_URL + "/security-question", {
      params: {
        user_id,
      },
    });
    const questions = response.data;
    return questions;
}

export async function fetchUserByEmail(email: string) {
  const response = await axios.get<User>(BASE_API_URL + "/user-email", {
    params: {
      email,
    },
  });
  const user = response.data;
  return user;
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
    const followingIds = new Set(
      followings?.map((user) => user.following_id) || []
    );
    // console.log(followingIds);

    const users = toFollows.map((user) => {
      return {
        ...user,
        is_followed: followingIds.has(user.id),
      };
    });

    return users;
}

export async function registerPremium(userId: string, identityNumber: string, facePicture: string, reason: string) {
  await axios.post(BASE_API_URL + "/register-premium", {
    userId,
    identityNumber,
    reason,
    facePicture,
  });
}