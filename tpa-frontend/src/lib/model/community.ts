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
  category: string;
  banner: string;
  rule: string;
};

const BASE_API_URL = "http://localhost:8080/api/threads";

export async function fetchCommunity() {
  const response = await axios.get<Community[]>(BASE_API_URL + "/community");
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