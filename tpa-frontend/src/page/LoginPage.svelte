<script lang="ts">
  let username = '';
  let password = '';
  let loading = false;
  
  async function handleLogin() {
    loading = true;
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      console.log('Logging in with:', { username, password });
      
      // Here you would handle the actual authentication
    } catch (error) {
      console.error('Login error:', error);
    } finally {
      loading = false;
    }
  }
  
  function handleGoogleLogin() {
    console.log('Google login clicked');
    // Implement Google OAuth login
  }
</script>

<div class="login-page">
  <div class="login-container">
    <div class="close-button">
      <button aria-label="Close">✕</button>
    </div>
    
    <div class="logo-container">
      <div class="logo">AY</div>
    </div>
    
    <h1>Sign in to AY</h1>
    
    <button class="google-button" on:click={handleGoogleLogin}>
      <svg viewBox="0 0 24 24" width="18" height="18">
        <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
        <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
        <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
        <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
      </svg>
      Sign in with Google
    </button>
    
    <div class="divider">
      <span>or</span>
    </div>
    
    <form on:submit|preventDefault={handleLogin}>
      <div class="input-group">
        <input 
          type="text" 
          placeholder="Phone, email, or username" 
          bind:value={username} 
          required
        />
      </div>
      
      <div class="input-group">
        <input 
          type="password" 
          placeholder="Password" 
          bind:value={password} 
          required
        />
      </div>
      
      <button 
        type="submit" 
        class="sign-in-button" 
        disabled={loading || !username || !password}
      >
        {#if loading}
          <span class="spinner"></span>
        {:else}
          Sign in
        {/if}
      </button>
    </form>
    
    <button class="forgot-password">
      Forgot password?
    </button>
    
    <div class="signup-prompt">
      Don't have an account? <a href="/signup">Sign up</a>
    </div>
  </div>
</div>

<style>
  /* Global styles */
  :global(html),
  :global(body) {
    margin: 0;
    padding: 0;
    overflow-x: hidden;
    width: 100%;
    height: 100%;
    background-color: #15202b;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }

  /* Login page wrapper */
  .login-page {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    width: 100%;
    padding: 20px;
    box-sizing: border-box;
  }
  
  /* Login container */
  .login-container {
    width: 100%;
    max-width: 300px;
    padding: 30px 20px;
    background-color: #000;
    color: #fff;
    display: flex;
    flex-direction: column;
    align-items: center;
    border-radius: 16px;
    position: relative;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
  }
  
  /* Close button */
  .close-button {
    position: absolute;
    top: 15px;
    left: 15px;
  }
  
  .close-button button {
    background: none;
    border: none;
    color: #fff;
    font-size: 16px;
    cursor: pointer;
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  /* Logo */
  .logo-container {
    margin-top: 20px;
    margin-bottom: 10px;
  }
  
  .logo {
    font-size: 28px;
    font-weight: bold;
  }
  
  /* Heading */
  h1 {
    font-size: 22px;
    margin-bottom: 24px;
    text-align: center;
    font-weight: 700;
  }
  
  /* Google button */
  .google-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    padding: 10px 0;
    border-radius: 50px;
    border: 1px solid #333;
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
  
  /* Divider */
  .divider {
    width: 100%;
    text-align: center;
    border-bottom: 1px solid #333;
    line-height: 0.1em;
    margin: 20px 0;
  }
  
  .divider span {
    background: #000;
    padding: 0 10px;
    color: #888;
    font-size: 14px;
  }
  
  /* Form */
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
    border: 1px solid #333;
    background-color: transparent;
    color: #fff;
    font-size: 14px;
    box-sizing: border-box;
  }
  
  input:focus {
    outline: none;
    border-color: #1d9bf0;
  }
  
  input::placeholder {
    color: #777;
  }
  
  /* Sign in button */
  .sign-in-button {
    width: 100%;
    padding: 12px;
    border-radius: 50px;
    border: none;
    background-color: #fff;
    color: #000;
    font-weight: bold;
    cursor: pointer;
    position: relative;
    margin-bottom: 15px;
    font-size: 15px;
    height: 44px;
  }
  
  .sign-in-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  /* Forgot password button */
  .forgot-password {
    width: 100%;
    padding: 12px;
    border-radius: 50px;
    border: 1px solid #333;
    background-color: transparent;
    color: #fff;
    font-weight: 500;
    cursor: pointer;
    margin-bottom: 20px;
    font-size: 15px;
    transition: background-color 0.2s;
  }
  
  .forgot-password:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
  
  /* Sign up prompt */
  .signup-prompt {
    font-size: 14px;
    color: #888;
    margin-top: 10px;
  }
  
  .signup-prompt a {
    color: #1d9bf0;
    text-decoration: none;
  }
  
  .signup-prompt a:hover {
    text-decoration: underline;
  }
  
  /* Loading spinner */
  .spinner {
    display: inline-block;
    width: 16px;
    height: 16px;
    border: 2px solid rgba(0, 0, 0, 0.3);
    border-radius: 50%;
    border-top-color: #000;
    animation: spin 1s ease-in-out infinite;
    position: absolute;
    left: calc(50% - 8px);
    top: calc(50% - 8px);
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  /* Responsive adjustments */
  @media (max-width: 400px) {
    .login-container {
      padding: 25px 15px;
    }
  }
</style>