<script>
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";


    onMount(async() => {
        const accessToken = localStorage.getItem("access-token");
        const res = await fetch("http://localhost:8080/api/auth/check", {
            headers: {
                "Authorization": `Bearer ${accessToken}`
            }
        });

        if(res.status === 401){
            // change to check refresh token
            console.log("not authenticated, redirecting back to login page")
            navigate("/login")
        }

        const refreshToken = localStorage.getItem("refresh-token");
        const response = await fetch("http://localhost:8080/api/auth/refresh", {
            
        })
    })

</script>

<h1>This is Home Page</h1>

<style>

</style>