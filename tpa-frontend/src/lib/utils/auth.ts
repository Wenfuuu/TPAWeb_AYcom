import axios from "axios";

const BASE_API_URL = "http://localhost:8080/api/auth";

type LoginResponse = {
    access_token: string
    refresh_token: string
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
            // meant user are logged in or authenticated, so redirect to home page
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
    navigate("/home");
    } else {
        // toast
    }
}