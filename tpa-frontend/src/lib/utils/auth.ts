export async function checkAuth(navigate: (path: string) => void) {
    const accessToken = localStorage.getItem("access-token");
    // console.log(accessToken)
    const res = await fetch("http://localhost:8080/api/auth/check", {
        headers: {
            "Authorization": `Bearer ${accessToken}`
        }
    });

    if(!res.ok){
        // change to check refresh token
        const refreshToken = localStorage.getItem("refresh-token");
        // console.log(refreshToken);
        const response = await fetch("http://localhost:8080/api/auth/refresh", {
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
    // console.log(accessToken)
    const res = await fetch("http://localhost:8080/api/auth/admin", {
        headers: {
            "Authorization": `Bearer ${accessToken}`
        }
    });

    if(!res.ok){
        // change this later to home page if not admin (not authorized)
        navigate("/");
    }
}