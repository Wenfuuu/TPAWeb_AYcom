import axios from "axios";
import type { User } from "../model/user";

const BASE_API_URL = "http://localhost:8080/api/auth";

type LoginResponse = {
    access_token: string
    refresh_token: string
    user: User
}

export function storeUserInfo(user: User): void {
    localStorage.setItem("user-info", JSON.stringify(user));
}

export function getUserInfo(): User | null {
    const userInfo = localStorage.getItem("user-info");
    return userInfo ? JSON.parse(userInfo) : null;
}

// get current user from API
export async function getAuthUser(): Promise<User | null> {
    const accessToken = localStorage.getItem("access-token");
    if (!accessToken) return null;

    try {
        const res = await axios.get(`${BASE_API_URL}/get-auth-user`, {
            headers: {
                "Authorization": `Bearer ${accessToken}`
            }
        });

        if (res.status === 200) {
            const user = res.data.user as User;
            // console.log(user.id);
            storeUserInfo(user);
            return user;
        }
    } catch (error) {
        console.error("Failed to get user:", error);
    }

    return null;
}

export async function checkGuest(navigate: (path: string) => void) {
    const accessToken = localStorage.getItem("access-token");
    // console.log(accessToken)
    const res = await fetch(BASE_API_URL + "/guest", {
        headers: {
            "Authorization": `Bearer ${accessToken}`
        }
    });

    // guest access token
    if(res.ok) {
        // change to check refresh token
        const refreshToken = localStorage.getItem("refresh-token");
        // console.log(refreshToken);
        const response = await fetch(BASE_API_URL + "/refresh", {
            headers: {
                "Refresh-token": `${refreshToken}`
            }
        });

        if(response.ok){
            // means user are logged in or authenticated, so redirect to home page
            navigate("/home");
        }
    }else{
        navigate("/home");
    }
}

export async function checkAuth(navigate: (path: string) => void) {
    const accessToken = localStorage.getItem("access-token");
    // console.log(accessToken)
    const res = await fetch(BASE_API_URL + "/check", {
        headers: {
            "Authorization": `Bearer ${accessToken}`
        }
    });

    if(!res.ok){
        // change to check refresh token
        const refreshToken = localStorage.getItem("refresh-token");
        // console.log(refreshToken);
        const response = await fetch(BASE_API_URL + "/refresh", {
            headers: {
                "Refresh-token": `${refreshToken}`
            }
        });

        if(!response.ok){
            navigate("/login");
        }else{
            // store new access & refresh token in local storage
            const data = await response.json();
            localStorage.setItem("access-token", data.access_token);
            localStorage.setItem("refresh-token", data.refresh_token);
        }
    }
}

export async function checkAdmin(navigate: (path: string) => void) {
    const accessToken = localStorage.getItem("access-token");

    await axios.get(BASE_API_URL + "/admin", {
        headers: {
            Authorization: `Bearer ${accessToken}`,
        },
    }).catch(() => {
        // change this later to home page if not admin
        navigate("/home");
    });
}

export async function login(email: string, password: string, navigate: (path: string) => void) {
    const res = await axios.post(BASE_API_URL + "/login", { email, password });

    if (res.status >= 200 && res.status < 300) {
        // Store tokens if present in response
        const data: LoginResponse = res.data;
        localStorage.setItem("access-token", data.access_token);
        localStorage.setItem("refresh-token", data.refresh_token);

        // Store user info
        if (data.user) {
            storeUserInfo(data.user);
        }

        navigate("/home");
    } else {
        // toast
    }
}

export function logout(navigate: (path: string) => void) {
    localStorage.removeItem("access-token");
    localStorage.removeItem("refresh-token");
    localStorage.removeItem("user-info");
    navigate("/login");
}