<script lang="ts">
  import { onMount } from "svelte";
  import { checkGuest } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import { toastStore } from "../lib/stores/toast";
  import { fetchSecurityQuestions, fetchUserByEmail, resetPassword, type SecurityQuestion, type User } from "../lib/model/user";
  import axios from "axios";
//   import { verifyEmail, getSecurityQuestion, verifySecurityAnswer, resetPassword } from "../lib/api/auth";

  // Recovery flow steps as string constants instead of enum
  const STEP_EMAIL_ENTRY = "email_entry";
  const STEP_SECURITY_QUESTION = "security_question";
  const STEP_NEW_PASSWORD = "new_password";
  const STEP_SUCCESS = "success";

  // Form data
  let email = "";
//   let securityAnswer = "";
  let newPassword = "";
  let confirmPassword = "";

  // UI state
  let currentStep = STEP_EMAIL_ENTRY;
  let isLoading = false;
//   let securityQuestion = "";
  let isDarkMode = true;

  // Validation state
  let emailError = "";
  let answerError = "";
  let passwordError = "";
  let confirmPasswordError = "";

  //to be fetched
  let user: User | null = null;
  let questions: SecurityQuestion;
  let securityQuestion = "";
  let securityAnswer = "";
  let oldPassword = "";

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }

  async function handleEmailSubmit() {
    if (!validateEmail()) return;

    isLoading = true;

    try {
        user = await fetchUserByEmail(email);
        console.log(user);

        questions = await fetchSecurityQuestions(user.id);
        console.log(questions);

        // get security details
        securityQuestion = questions.security_question;

        // Move to security question step
        currentStep = STEP_SECURITY_QUESTION;
    } catch (error) {
        if (axios.isAxiosError(error)) {
            const message = error.response?.data?.error || error.message;
            toastStore.error(message);
        } else if (error instanceof Error) {
            toastStore.error(error.message);
        }
    } finally {
        isLoading = false;
    }
  }

  async function handleSecurityAnswerSubmit() {
    if (!validateSecurityAnswer()) return;

    isLoading = true;

    try {
      const isCorrect = securityAnswer.toLowerCase() === questions.security_answer.toLowerCase();

      if (isCorrect) {
        // Answer is correct, move to password reset
        currentStep = STEP_NEW_PASSWORD;
      } else {
        answerError = "Incorrect answer. Please try again.";
      }

    } catch (error) {
        if (axios.isAxiosError(error)) {
            const message = error.response?.data?.error || error.message;
            toastStore.error(message);
        } else if (error instanceof Error) {
            toastStore.error(error.message);
        }
    } finally {
        isLoading = false;
    }
  }

  async function handlePasswordReset() {
    if (!validatePassword() || !user) return;

    isLoading = true;

    try {
      await resetPassword(user.id, newPassword);
      // Password reset successful
      currentStep = STEP_SUCCESS;
    } catch (error) {
        if (axios.isAxiosError(error)) {
            const message = error.response?.data?.error || error.message;
            toastStore.error(message);
        } else if (error instanceof Error) {
            toastStore.error(error.message);
        }
    } finally {
      isLoading = false;
    }
  }

  function validateEmail() {
    emailError = "";

    if (!email.trim()) {
      emailError = "Email is required";
      return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      emailError = "Please enter a valid email address";
      return false;
    }

    return true;
  }

  function validateSecurityAnswer() {
    answerError = "";

    if (!securityAnswer.trim()) {
      answerError = "Please provide an answer to the security question";
      return false;
    }

    return true;
  }

  function validatePassword() {
    passwordError = "";
    confirmPasswordError = "";

    if (!newPassword) {
      passwordError = "Password is required";
      return false;
    }

    if (newPassword.length < 8) {
      passwordError = "Password must be at least 8 characters long";
      return false;
    }

    if (!confirmPassword) {
      confirmPasswordError = "Please confirm your password";
      return false;
    }

    if (newPassword !== confirmPassword) {
      confirmPasswordError = "Passwords do not match";
      return false;
    }

    return true;
  }

  function goToLogin() {
    navigate("/login");
  }

  function goToLanding() {
    navigate("/");
  }

  onMount(async() => {
    await checkGuest(navigate);

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
      // Check user's system preference if no saved theme
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute("data-theme", prefersDark ? "dark" : "light");
    }
  });
</script>

<main class="forgotten-page">
  <div class="top-controls">
    <button
      type="button"
      class="theme-toggle"
      aria-label="Toggle theme"
      on:click={toggleTheme}
    >
      {#if isDarkMode}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
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
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
        </svg>
      {/if}
    </button>
  </div>

  <div class="content">
    <div class="form-container">
      <h1>Account Recovery</h1>

      {#if currentStep === STEP_EMAIL_ENTRY}
        <form on:submit|preventDefault={handleEmailSubmit}>
          <div class="form-group">
            <label for="email">Email</label>
            <input
              type="email"
              id="email"
              bind:value={email}
              class:error={emailError}
              placeholder="Enter your email address"
              disabled={isLoading}
            />
            {#if emailError}
              <span class="error-message">{emailError}</span>
            {/if}
          </div>

          <button
            type="submit"
            class="primary-button"
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="spinner"></span>
              Verifying...
            {:else}
              Continue
            {/if}
          </button>
        </form>
      {:else if currentStep === STEP_SECURITY_QUESTION}
        <form on:submit|preventDefault={handleSecurityAnswerSubmit}>
          <div class="form-group">
            <p>Security Question</p>
            <div class="security-question">{securityQuestion}</div>
          </div>

          <div class="form-group">
            <label for="security-answer">Your Answer</label>
            <input
              type="text"
              id="security-answer"
              bind:value={securityAnswer}
              class:error={answerError}
              placeholder="Enter your answer"
              disabled={isLoading}
            />
            {#if answerError}
              <span class="error-message">{answerError}</span>
            {/if}
          </div>

          <button
            type="submit"
            class="primary-button"
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="spinner"></span>
              Verifying...
            {:else}
              Verify Answer
            {/if}
          </button>

          <button
            type="button"
            class="text-button"
            on:click={() => currentStep = STEP_EMAIL_ENTRY}
            disabled={isLoading}
          >
            Back to Email
          </button>
        </form>
      {:else if currentStep === STEP_NEW_PASSWORD}
        <form on:submit|preventDefault={handlePasswordReset}>
          <div class="form-group">
            <label for="new-password">New Password</label>
            <input
              type="password"
              id="new-password"
              bind:value={newPassword}
              class:error={passwordError}
              placeholder="Enter new password"
              disabled={isLoading}
            />
            {#if passwordError}
              <span class="error-message">{passwordError}</span>
            {/if}
            <div class="password-hint">Password must be at least 8 characters long</div>
          </div>

          <div class="form-group">
            <label for="confirm-password">Confirm Password</label>
            <input
              type="password"
              id="confirm-password"
              bind:value={confirmPassword}
              class:error={confirmPasswordError}
              placeholder="Confirm new password"
              disabled={isLoading}
            />
            {#if confirmPasswordError}
              <span class="error-message">{confirmPasswordError}</span>
            {/if}
          </div>

          <button
            type="submit"
            class="primary-button"
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="spinner"></span>
              Resetting...
            {:else}
              Reset Password
            {/if}
          </button>

          <button
            type="button"
            class="text-button"
            on:click={() => currentStep = STEP_SECURITY_QUESTION}
            disabled={isLoading}
          >
            Back to Security Question
          </button>
        </form>
      {:else if currentStep === STEP_SUCCESS}
        <div class="success-container">
          <div class="success-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="64"
              height="64"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
          </div>
          <h2>Password Reset Successful</h2>
          <p>
            Your password has been reset successfully. You can now log in with your new password.
          </p>
          <button
            type="button"
            class="primary-button"
            on:click={goToLogin}
          >
            Go to Login
          </button>
        </div>
      {/if}

      <div class="form-footer">
        <p>
          <a href="/" on:click|preventDefault={goToLanding}>Return to Landing Page</a>
        </p>
      </div>
    </div>
  </div>
</main>

<style>
  .forgotten-page {
    width: 100%;
    height: 100vh;
    position: relative;
    overflow-y: auto;
    overflow-x: hidden;
    background-color: var(--background);
    color: var(--text-primary);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .top-controls {
    position: absolute;
    top: 20px;
    right: 20px;
    z-index: 10;
  }

  .theme-toggle {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .theme-toggle:hover {
    background-color: var(--hover-background);
  }

  .content {
    width: 100%;
    max-width: 400px;
    padding: 0 20px;
    box-sizing: border-box;
  }

  .form-container {
    background-color: var(--card-background);
    border-radius: 16px;
    padding: 24px;
    border: 1px solid var(--border-color);
    margin-bottom: 24px;
  }

  h1 {
    font-size: 24px;
    font-weight: 700;
    margin: 0 0 24px 0;
    text-align: center;
  }

  .form-group {
    margin-bottom: 20px;
  }

  label {
    display: block;
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
  }

  input {
    width: 100%;
    background-color: var(--hover-background);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 12px;
    color: var(--text-primary);
    font-size: 15px;
    transition: border-color 0.2s;
    box-sizing: border-box;
  }

  input:focus {
    border-color: var(--button-primary);
    outline: none;
  }

  input.error {
    border-color: #ff4d4f;
  }

  .error-message {
    display: block;
    color: #ff4d4f;
    font-size: 14px;
    margin-top: 4px;
  }

  .password-hint {
    font-size: 12px;
    color: var(--text-tertiary);
    margin-top: 4px;
  }

  .security-question {
    background-color: var(--hover-background);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 12px;
    font-size: 15px;
    margin-bottom: 8px;
  }

  .primary-button {
    width: 100%;
    background-color: var(--button-primary);
    border: none;
    color: white;
    padding: 12px;
    border-radius: 9999px;
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-bottom: 16px;
  }

  .primary-button:hover:not(:disabled) {
    background-color: var(--button-primary-hover);
  }

  .primary-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .text-button {
    width: 100%;
    background: none;
    border: none;
    color: var(--button-primary);
    padding: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: color 0.2s;
    margin-bottom: 8px;
  }

  .text-button:hover:not(:disabled) {
    text-decoration: underline;
  }

  .text-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .spinner {
    display: inline-block;
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 1s ease-in-out infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .form-footer {
    text-align: center;
    margin-top: 16px;
    font-size: 14px;
  }

  .form-footer a {
    color: var(--button-primary);
    text-decoration: none;
  }

  .form-footer a:hover {
    text-decoration: underline;
  }

  .success-container {
    text-align: center;
  }

  .success-icon {
    display: flex;
    justify-content: center;
    margin-bottom: 16px;
  }

  .success-icon svg {
    color: #52c41a;
  }

  .success-container h2 {
    font-size: 20px;
    font-weight: 700;
    margin: 0 0 16px 0;
  }

  .success-container p {
    font-size: 15px;
    line-height: 1.5;
    margin-bottom: 24px;
  }

  @media (max-width: 480px) {
    .form-container {
      padding: 20px;
    }
  }
</style>