import axios from "axios";

export type CommunityCategory = {
  id: string;
  name: string;
};

export type Community = {
  id: string;
  name: string;
  description: string;
  logo: string;
  category_id: string;
  banner: string;
  rule: string;
  category: string;
};

const BASE_API_URL = "http://localhost:8080/api/threads";

let communities: Community[] = [];

export async function fetchCommunityWithCategory(page: number, limit: number, search: string) {
    const paginatedCommunity = await fetchCommunityPagination(page, limit, search);
    const categories = await fetchCommunityCategory();

    if(!paginatedCommunity) return null;

    communities = paginatedCommunity.map(community => {
        const c = categories.find(c => c.id === community.category_id);
        const category: string = c ? c.name : "Unknown";

        return {
          ...community,
          category,
        };
    });

    return communities;
}

export async function fetchPendingWithCategory(userid: string, page: number, limit: number, search: string) {
    const paginatedCommunity = await fetchPendingCommunity(userid, page, limit, search);
    const categories = await fetchCommunityCategory();

    if(!paginatedCommunity) return null;

    communities = paginatedCommunity.map(community => {
        const c = categories.find(c => c.id === community.category_id);
        const category: string = c ? c.name : "Unknown";

        return {
          ...community,
          category,
        };
    });

    return communities;
}

export async function fetchCommunity() {
  const response = await axios.get<Community[]>(BASE_API_URL + "/community");
  const communities = response.data;
  return communities;
}

async function fetchCommunityPagination(page: number, limit: number, search: string) {
    const response = await axios.get<Community[]>(BASE_API_URL + "/community-pagination", {
        params: {
            page,
            limit,
            search
        }
    });
    const communities = response.data;
    return communities;
}

async function fetchPendingCommunity(user_id: string, page: number, limit: number, search: string) {
  const response = await axios.get<Community[]>(BASE_API_URL + "/pending-community", {
        params: {
          user_id,
          page,
          limit,
          search
        }
    });
    const communities = response.data;
    return communities;
}

export async function fetchCommunityCategory() {
  const response = await axios.get<CommunityCategory[]>(
    BASE_API_URL + "/community-category"
  );
  const categories = response.data;
  return categories;
}