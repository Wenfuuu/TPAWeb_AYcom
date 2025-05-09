<script module lang="ts">
  declare global {
    interface Window {
      recaptchaCallback: (token: string) => void;
    }
  }
</script>

<script lang="ts">
  import { navigate } from "svelte-routing";
  import { checkGuest, login } from "../lib/utils/auth";
  import { onMount } from "svelte";
  import { toastStore } from "../lib/stores/toast";

  let email = "";
  let password = "";
  let loading = false;
  let recaptchaVerified = false;
  let isDarkMode = true;

  async function handleLogin() {
    loading = true;

    try {
      await login(email, password, navigate);
      toastStore.success("Login success!");
    } catch (error: unknown) {
      if (error instanceof Error) {
        toastStore.error(error.message);
      } else {
        toastStore.error("An unexpected error occurred!");
      }
    } finally {
      loading = false;
    }
  }

  function handleGoogleLogin() {
    // Implement Google OAuth login
  }

  function navigateToForgottenPage() {
    navigate("/forgotten");
  }

  function onRecaptchaCallback(token: string): void {
    fetch("http://localhost:8080/verify-recaptcha", {
      method: "POST",
      // headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ token })
    })
      .then(res => res.json())
      .then(data => {
        if (data.success) {
          // console.log("success recaptcha");
          recaptchaVerified = true;
        }
      });
  }

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }

  onMount(async () => {
    // console.log("Component mounted, reCAPTCHA would be initialized here");
    // Assign the global callback used by reCAPTCHA v2
    window.recaptchaCallback = onRecaptchaCallback;
    // Dynamically load the reCAPTCHA script if not already on the page
    if (!document.querySelector("script[src=\"https://www.google.com/recaptcha/api.js\"]")) {
      const script = document.createElement("script");
      script.src = "https://www.google.com/recaptcha/api.js";
      script.async = true;
      script.defer = true;
      document.body.appendChild(script);
    }

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute("data-theme", prefersDark ? "dark" : "light");
    }

    // add check guest auth here
    await checkGuest(navigate);
  });
</script>

<div class="login-page">
  <div class="login-container">
    <div class="close-button">
      <button type="button" aria-label="Close">✕</button>
    </div>

    <div class="theme-toggle">
      <button type="button" on:click={toggleTheme} aria-label="Toggle theme">
        {#if isDarkMode}
          <!-- Sun icon for light mode -->
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="5"></circle>
            <line x1="12" y1="1" x2="12" y2="3"></line>
            <line x1="12" y1="21" x2="12" y2="23"></line>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
            <line x1="1" y1="12" x2="3" y2="12"></line>
            <line x1="21" y1="12" x2="23" y2="12"></line>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
          </svg>
        {:else}
          <!-- Moon icon for dark mode -->
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
          </svg>
        {/if}
      </button>
    </div>

    <div class="logo-container">
      {#if isDarkMode}
        <img src="/src/assets/AY-dark.jpg" alt="Logo" class="logo-img">
      {:else}
        <img src="/src/assets/AY-light.jpg" alt="Logo" class="logo-img">
      {/if}
    </div>

    <h1>Sign in to AY</h1>

    <button type="button" class="google-button" on:click={handleGoogleLogin}>
      <svg viewBox="0 0 24 24" width="18" height="18">
        <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
        <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
        <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
        <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
      </svg>
      Sign in with Google
    </button>

    <div class="divider">
      <span>or</span>
    </div>

    <form>
      <div class="input-group">
        <input
          type="text"
          placeholder="Email"
          bind:value={email}
        />
      </div>

      <div class="input-group">
        <input
          type="password"
          placeholder="Password"
          bind:value={password}
        />
      </div>

      <!-- recaptcha here -->
      <div class="recaptcha">
          <div
            id="recaptcha-container"
            class="g-recaptcha"
            data-sitekey="6Le5ySYrAAAAADAVGZm_Gug20l8ge7D0IQgwg6f3"
            data-callback="recaptchaCallback">
          </div>
      </div>

      <button
        type="button"
        class="sign-in-button"
        disabled={loading || !email || !password || !recaptchaVerified}
        on:click={handleLogin}
      >
        {#if loading}
          <span class="spinner"></span>
        {:else}
          Sign in
        {/if}
      </button>
    </form>

    <button class="forgot-password"
    type="button"
    on:click={navigateToForgottenPage}>
      Forgot password?
    </button>

    <div class="signup-prompt">
      Don't have an account? <a href="/register">Sign up</a>
    </div>
  </div>
</div>

<style>
  :global(body) {
    background-color: var(--background);
    color: var(--text-primary);
  }

  .login-page {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    width: 100%;
    padding: 20px;
    box-sizing: border-box;
  }

  .login-container {
    width: 100%;
    max-width: 300px;
    padding: 30px 20px;
    background-color: var(--card-background);
    color: var(--text-primary);
    display: flex;
    flex-direction: column;
    align-items: center;
    border-radius: 16px;
    position: relative;
    box-shadow: 0 8px 30px var(--shadow-color);
    transition: background-color 0.3s ease, color 0.3s ease, box-shadow 0.3s ease;
  }

  .close-button {
    position: absolute;
    top: 15px;
    left: 15px;
  }

  .close-button button {
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 16px;
    cursor: pointer;
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .theme-toggle {
    position: absolute;
    top: 15px;
    right: 15px;
  }

  .theme-toggle button {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: background-color 0.2s;
  }

  .theme-toggle button:hover {
    background-color: var(--button-secondary-hover);
  }

  .logo-container {
    margin-top: 20px;
    margin-bottom: 10px;
  }

  h1 {
    font-size: 22px;
    margin-bottom: 24px;
    text-align: center;
    font-weight: 700;
    color: var(--text-primary);
  }

  .logo-img {
    max-width: 50px;
  }

  .google-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    padding: 10px 0;
    border-radius: 50px;
    border: 1px solid var(--border-color);
    background-color: #fff;
    color: #000;
    font-weight: 500;
    margin-bottom: 15px;
    cursor: pointer;
    transition: background-color 0.2s;
    font-size: 15px;
  }

  .google-button:hover {
    background-color: #f5f5f5;
  }

  .google-button svg {
    margin-right: 8px;
  }

  .divider {
    width: 100%;
    text-align: center;
    border-bottom: 1px solid var(--divider-color);
    line-height: 0.1em;
    margin: 20px 0;
  }

  .divider span {
    background: var(--divider-text-background);
    padding: 0 10px;
    color: var(--text-tertiary);
    font-size: 14px;
    transition: background-color 0.3s ease;
  }

  form {
    width: 100%;
  }

  .input-group {
    width: 100%;
    margin-bottom: 15px;
  }

  input {
    width: 100%;
    padding: 12px;
    border-radius: 4px;
    border: 1px solid var(--input-border);
    background-color: var(--input-background);
    color: var(--text-primary);
    font-size: 14px;
    box-sizing: border-box;
    transition: border-color 0.3s ease, background-color 0.3s ease, color 0.3s ease;
  }

  input:focus {
    outline: none;
    border-color: var(--button-primary);
  }

  input::placeholder {
    color: var(--text-tertiary);
  }

  .recaptcha {
    margin-bottom: 24px;
  }

  .sign-in-button {
    width: 100%;
    padding: 12px;
    border-radius: 50px;
    border: none;
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    font-weight: bold;
    cursor: pointer;
    position: relative;
    margin-bottom: 15px;
    font-size: 15px;
    height: 44px;
    transition: background-color 0.2s;
  }

  .sign-in-button:hover:not(:disabled) {
    background-color: var(--button-primary-hover);
  }

  .sign-in-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .forgot-password {
    width: 100%;
    padding: 12px;
    border-radius: 50px;
    border: 1px solid var(--button-secondary-border);
    background-color: var(--button-secondary-bg);
    color: var(--button-secondary-text);
    font-weight: 500;
    cursor: pointer;
    margin-bottom: 20px;
    font-size: 15px;
    transition: background-color 0.2s;
  }

  .forgot-password:hover {
    background-color: var(--button-secondary-hover);
  }

  .signup-prompt {
    font-size: 14px;
    color: var(--text-tertiary);
    margin-top: 10px;
  }

  .signup-prompt a {
    color: var(--button-primary);
    text-decoration: none;
  }

  .signup-prompt a:hover {
    text-decoration: underline;
  }

  .spinner {
    display: inline-block;
    width: 16px;
    height: 16px;
    border: 2px solid var(--spinner-border);
    border-radius: 50%;
    border-top-color: var(--spinner-border-top);
    animation: spin 1s ease-in-out infinite;
    position: absolute;
    left: calc(50% - 8px);
    top: calc(50% - 8px);
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  @media (max-width: 400px) {
    .login-container {
      padding: 25px 15px;
    }
  }
</style>