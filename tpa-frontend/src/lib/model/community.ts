import axios from "axios";
import { fetchUser } from "./user";

export type CommunityCategory = {
  id: string;
  name: string;
};

export type CommunityMember = {
  id: string;
  user_id: string;
  community_id: string;
  role: string;
  name: string;
  username: string;
  profile_picture: string;
  is_premium: boolean;
};

export type JoinCommunityRequest = {
  id: string;
  user_id: string;
  community_id: string;
  created_at: string;
  name: string;
  username: string;
};

export type CreateCommunityRequest = {
  id: string;
  user_id: string;
  name: string;
  description: string;
  icon_url: string;
  banner_url: string;
  category_id: string;
  rule: string;
  created_at: string;
  creator_name: string;
};

export type CommunityMedia = {
  id: string;
  thread_id: string;
  file: string;
}

export type Community = {
  id: string;
  name: string;
  description: string;
  logo: string;
  category_id: string;
  banner: string;
  rule: string;
  category: string;
  creator_id: string;
  created_at: string;
  creator_name: string;
  is_joined: boolean;
  is_pending: boolean;
};

const BASE_API_URL = "http://localhost:8080/api/threads";

let communities: Community[] = [];

export async function requestJoinCommunity(userId: string, communityId: string) {
  await axios.post(BASE_API_URL + "/request-join-community", {
    userId,
    communityId,
  });
}

export async function acceptCreateCommunity(request_id: string) {
  await axios.post(BASE_API_URL + "/accept-create-community", {
    request_id
  });
}

export async function rejectCreateCommunity(request_id: string) {
  await axios.post(BASE_API_URL + "/reject-create-community", {
    request_id
  });
}

export async function fetchCreateCommunity() {
  const response = await axios.get(BASE_API_URL + "/request-create-community");
  // return response.data;

  const allRequests: CreateCommunityRequest[] = response.data;
  if(!allRequests) return [];
  const allUsers = await fetchUser();

  const requests = allRequests.map((request) => {
    const user = allUsers.find((u) => u.id === request.user_id);
    const creator_name = user ? user.name : "Unknown";

    return {
      ...request,
      creator_name
    };
  });

  return requests;
}

export async function requestCreateCommunity(userId: string, name: string, description: string,
  iconUrl: string, bannerUrl: string, categoryId: string, rule: string
) {
  // console.log("requesting community");
  await axios.post(BASE_API_URL + "/request-create-community", {
    userId,
    name,
    description,
    iconUrl,
    bannerUrl,
    categoryId,
    rule
  });
}

export async function fetchCommunityWithCategory(userid: string, page: number, limit: number, search: string, categoryid: string) {
  const paginatedCommunity = await fetchCommunityPagination(page, limit, search, categoryid);
  const categories = await fetchCommunityCategory();
  const allUsers = await fetchUser();
  let pendingCommunities = await fetchPendingCommunity(userid, 1, 1000, "", categoryid);
  let joinedCommunities = await fetchJoinedCommunity(userid, 1, 1000, "", categoryid);
  if(!pendingCommunities) pendingCommunities = [];
  if(!joinedCommunities) joinedCommunities = [];

  if (!paginatedCommunity) return null;

  const excludedIds = new Set([
    ...pendingCommunities.map((community) => community.id),
    ...joinedCommunities.map((community) => community.id),
  ]);

  const filteredCommunities = paginatedCommunity.filter(
    (community) => !excludedIds.has(community.id)
  );

  communities = filteredCommunities.map((community) => {
    const c = categories.find((c) => c.id === community.category_id);
    const category: string = c ? c.name : "Unknown";

    const user = allUsers.find((u) => u.id === community.creator_id);
    const creator_name = user ? user.name : "Unknown";

    return {
      ...community,
      category,
      creator_name,
    };
  });

  return communities;
}

export async function fetchPendingWithCategory(userid: string, page: number, limit: number, search: string, categoryid: string) {
    const paginatedCommunity = await fetchPendingCommunity(userid, page, limit, search, categoryid);
    const categories = await fetchCommunityCategory();
    const allUsers = await fetchUser();

    if(!paginatedCommunity) return null;

    communities = paginatedCommunity.map(community => {
        const c = categories.find(c => c.id === community.category_id);
        const category: string = c ? c.name : "Unknown";

        const user = allUsers.find((u) => u.id === community.creator_id);
        const creator_name = user ? user.name : "Unknown";

        return {
          ...community,
          category,
          creator_name
        };
    });

    return communities;
}

export async function fetchJoinedWithCategory(userid: string, page: number, limit: number, search: string, categoryid: string) {
    const paginatedCommunity = await fetchJoinedCommunity(userid, page, limit, search, categoryid);
    const categories = await fetchCommunityCategory();
    const allUsers = await fetchUser();

    if(!paginatedCommunity) return null;

    communities = paginatedCommunity.map(community => {
        const c = categories.find((c) => c.id === community.category_id);
        const category: string = c ? c.name : "Unknown";

        const user = allUsers.find((u) => u.id === community.creator_id);
        const creator_name = user ? user.name : "Unknown";

        return {
          ...community,
          category,
          creator_name,
        };
    });

    return communities;
}

export async function fetchCommunity() {
  const response = await axios.get<Community[]>(BASE_API_URL + "/community");
  const communities = response.data;
  return communities;
}

export async function fetchCommunityDetailWithCreator(user_id: string, community_id: string) {
  const community = await fetchCommunityDetail(community_id);
  const categories = await fetchCommunityCategory();
  const allUsers = await fetchUser();

  let pendingCommunities = await fetchPendingCommunity(user_id, 1, 1000, "", "all");
  let joinedCommunities = await fetchJoinedCommunity(user_id, 1, 1000, "", "all");
  if(!pendingCommunities) pendingCommunities = [];
  if(!joinedCommunities) joinedCommunities = [];

  const communityArray = Array.isArray(community) ? community : [community];

  communities = communityArray.map((community) => {
    const c = categories.find((c) => c.id === community.category_id);
    const category: string = c ? c.name : "Unknown";

    const user = allUsers.find((u) => u.id === community.creator_id);
    const creator_name = user ? user.name : "Unknown";

    const is_joined = joinedCommunities.some((joined) => joined.id === community.id);
    const is_pending = pendingCommunities.some((pending) => pending.id === community.id);

    return {
      ...community,
      category,
      creator_name,
      is_joined,
      is_pending,
    };
  });

  return communities[0];
}

export async function fetchCommunityDetail(community_id: string) {
  const response = await axios.get<Community>(BASE_API_URL + "/community-detail", {
    params: {
      community_id
    }
  });
  const communities = response.data;
  return communities;
}

export async function fetchJoinRequest(community_id: string) {
  let allRequests = await fetchJoinCommunityRequest(community_id);
  const allUsers = await fetchUser();

  if(!allRequests) allRequests = [];

  const requests = allRequests.map(request => {
    const user = allUsers.find((u) => u.id === request.user_id);
    const name = user ? user.name : "Unknown";
    const username = user ? user.username : "Unknown";

    return{
      ...request,
      name,
      username
    };
  });
  return requests;
}

async function fetchJoinCommunityRequest(community_id: string) {
  const response = await axios.get<JoinCommunityRequest[]>(BASE_API_URL + "/request-join-community", {
    params: {
      community_id
    }
  });
  const requests = response.data;
  return requests;
}

async function fetchCommunityPagination(page: number, limit: number, search: string, category_id: string) {
    const response = await axios.get<Community[]>(BASE_API_URL + "/community-pagination", {
        params: {
            page,
            limit,
            search,
            category_id
        }
    });
    const communities = response.data;
    return communities;
}

async function fetchPendingCommunity(user_id: string, page: number, limit: number, search: string, category_id: string) {
  const response = await axios.get<Community[]>(BASE_API_URL + "/pending-community", {
        params: {
          user_id,
          page,
          limit,
          search,
          category_id
        }
    });
    const communities = response.data;
    return communities;
}

async function fetchJoinedCommunity(user_id: string, page: number, limit: number, search: string, category_id: string) {
  const response = await axios.get<Community[]>(BASE_API_URL + "/joined-community", {
        params: {
          user_id,
          page,
          limit,
          search,
          category_id
        }
    });
    const communities = response.data;
    return communities;
}

export async function fetchTopMember(community_id: string) {
  let topMembers = await fetchTopCommunityMember(community_id);
  const allUsers = await fetchUser();

  if(!topMembers) topMembers = [];

  const members: CommunityMember[] = topMembers.map((member) => {
    const user = allUsers.find((u) => u.id === member.user_id);
    // console.log(user);
    const name = user ? user.name : "Unknown";
    const username = user ? user.username : "Unknown";
    const profile_picture = user ? user.profile_picture : "Unknown";
    const is_premium = user?.is_premium ? user.is_premium : false;
    const bio = user?.bio ? user.bio : "Unknown";

    return {
      ...member,
      name, username, profile_picture,
      is_premium, bio
    };
  });

  return members;
}

export async function fetchModerator(community_id: string) {
  let topMembers = await fetchCommunityModerator(community_id);
  const allUsers = await fetchUser();

  if(!topMembers) topMembers = [];

  const members: CommunityMember[] = topMembers.map((member) => {
    const user = allUsers.find((u) => u.id === member.user_id);
    // console.log(user);
    const name = user ? user.name : "Unknown";
    const username = user ? user.username : "Unknown";
    const profile_picture = user ? user.profile_picture : "Unknown";
    const is_premium = user?.is_premium ? user.is_premium : false;

    return {
      ...member,
      name,
      username,
      profile_picture,
      is_premium,
    };
  });

  return members;
}

async function fetchCommunityModerator(community_id: string) {
  const response = await axios.get<CommunityMember[]>(BASE_API_URL + "/community-member", {
    params: {
      community_id
    }
  });
  let members: CommunityMember[] = response.data;
  if(!members) members = [];
  const moderator = members.filter((m) => m.role !== "member");
  return moderator;
}

export async function fetchCommunityMember(community_id: string) {
  const response = await axios.get<CommunityMember[]>(BASE_API_URL + "/community-member", {
    params: {
      community_id
    }
  });
  const members = response.data;
  return members;
}

async function fetchTopCommunityMember(community_id: string) {
  const response = await axios.get<CommunityMember[]>(BASE_API_URL + "/top-community-member", {
    params: {
      community_id
    }
  });
  const communities = response.data;
  return communities;
}

export async function acceptJoinRequest(userId: string, communityId: string) {
  await axios.post(BASE_API_URL + "/accept-join-community", {
    userId,
    communityId
  });
}

export async function rejectJoinRequest(userId: string, communityId: string) {
  await axios.post(BASE_API_URL + "/reject-join-community", {
    userId,
    communityId,
  });
}

export async function createCommunityCategory(name: string) {
  await axios.post(BASE_API_URL + "/community-category", {
    name
  });
}

export async function updateCommCategory(category_id: string, name: string) {
  await axios.put(BASE_API_URL + "/community-category", {
    category_id,
    name
  });
}

export async function deleteCommCategory(category_id: string) {
  await axios.delete(BASE_API_URL + "/community-category", {
    params: {
      category_id
    }
  });
}

export async function fetchCommunityCategory() {
  const response = await axios.get<CommunityCategory[]>(
    BASE_API_URL + "/community-category"
  );
  const categories = response.data;
  if(!categories) return [];
  return categories;
}

export async function fetchCommunityMedia(page: number, community_id: string) {
  const response = await axios.get<CommunityMedia[]>(
    BASE_API_URL + "/community-media", {
      params: {
        page,
        community_id
      }
    }
  );
  const media = response.data;
  return media;
}