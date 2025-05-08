<script>
  import { onMount } from "svelte";
    import { navigate } from "svelte-routing";
  import { checkGuest } from "../lib/utils/auth";

    function register() {
        navigate("/register");
    }

    function login() {
        navigate("/login");
    }

    onMount(async() => {
        await checkGuest(navigate);
    });

</script>

<div class="landing-page">
    <div class="container">
        <div class="logo-section">
            <img src="/src/assets/AY-dark.jpg" alt="Logo" class="logo-img">
        </div>
        <div class="content-section">
        <h1>Happening now</h1>
        <h2>Join today.</h2>
        <p class="tagline">Connect, share, engage.</p>

        <div class="account-section">
            <p class="account-prompt">Don't have an account?</p>
            <button class="create-account-btn" type="button" onclick={register}>Create account</button>
            <p class="terms">
            By signing up, you agree to the <a href="/terms">Terms of Service</a> and <a href="/privacy">Privacy Policy</a>, including <a href="/cookies">Cookie Use</a>.
            </p>

            <div class="signin-section">
            <p class="account-prompt">Already have an account?</p>
            <button class="sign-in-btn" type="button" onclick={login}>Sign in</button>
            </div>
        </div>
        </div>
    </div>
</div>

<style>
    :global(html), :global(body) {
        margin: 0;
        padding: 0;
        overflow-x: hidden;
        width: 100%;
        height: 100%;
    }

    .landing-page {
        color-scheme: dark light;
        background-color: #000;
        color: #fff;
        min-height: 100vh;
        display: flex;
        align-items: center;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    }

    .container {
        width: 100%;
        max-width: 1200px;
        margin: 0 auto;
        padding: 2rem;
        display: flex;
        gap: 2rem;
    }

    .logo-section {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .logo {
        font-size: 10rem;
        font-weight: 900;
        line-height: 1;
    }

    .logo-img {
        max-width: 300px;
    }

    .content-section {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    h1 {
        font-size: 4rem;
        font-weight: 800;
        margin: 0 0 1.5rem 0;
    }

    h2 {
        font-size: 2rem;
        font-weight: 700;
        margin: 0 0 0.5rem 0;
    }

    .tagline {
        font-size: 1.25rem;
        margin: 0 0 2rem 0;
    }

    .account-section {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        max-width: 300px;
    }

    .account-prompt {
        margin: 0;
    }

    .create-account-btn {
        background-color: #1d9bf0;
        color: white;
        border: none;
        border-radius: 9999px;
        padding: 0.75rem 0;
        font-size: 1rem;
        font-weight: 700;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .create-account-btn:hover {
        background-color: #1a8cd8;
    }

    .terms {
        font-size: 0.75rem;
        color: #71767b;
        margin: 0;
    }

    .terms a {
        color: #1d9bf0;
        text-decoration: none;
    }

    .terms a:hover {
        text-decoration: underline;
    }

    .signin-section {
        margin-top: 1.5rem;
    }

    .sign-in-btn {
        background-color: transparent;
        color: #1d9bf0;
        border: 1px solid #536471;
        border-radius: 9999px;
        padding: 0.75rem 0;
        width: 100%;
        font-size: 1rem;
        font-weight: 700;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .sign-in-btn:hover {
        background-color: rgba(29, 155, 240, 0.1);
    }

    /* Responsive styles */
    @media (max-width: 768px) {
        .container {
        flex-direction: column;
        padding: 1rem;
        }

        .logo {
        font-size: 5rem;
        }

        h1 {
        font-size: 2.5rem;
        }

        h2 {
        font-size: 1.5rem;
        }
    }
</style>