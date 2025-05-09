<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import { registerPremium, type User } from "../lib/model/user";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { uploadImage } from "../lib/utils/supabaseStorage";
  import axios from "axios";
  import { toastStore } from "../lib/stores/toast";

  let user: User | null = null;
  let idNumber = "";
  let verifyReason = "";
  let faceImage: File | null = null;
  let faceImagePreview: string | null = null;
  let isSubmitting = false;
  let isSubmitted = false;
  let isDarkMode = true;

  let idNumberValid = true;
  let imageValid = true;
  let reasonValid = true;

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }

  function handleImageUpload(event: Event) {
    const target = event.target as HTMLInputElement;
    const files = target.files;

    if (!files || files.length === 0) {
      faceImage = null;
      faceImagePreview = null;
      return;
    }

    const file = files[0];
    faceImage = file;

    const reader = new FileReader();
    reader.onload = (e: ProgressEvent<FileReader>) => {
      const result = e.target?.result;
      if (typeof result === "string") {
        faceImagePreview = result;
      }
    };

    reader.readAsDataURL(file);
  }

  function validateForm(): boolean {
    idNumberValid = idNumber.trim().length >= 5;
    imageValid = faceImage !== null;
    reasonValid = verifyReason.trim().length >= 10;

    return idNumberValid && imageValid && reasonValid;
  }

  async function submitVerification() {
    if (!validateForm() || !faceImage || !user) {
      return;
    }

    isSubmitting = true;
    const supaURL = await uploadImage(faceImage);
    if(!supaURL) return;

    try{
      await registerPremium(user.id, idNumber, supaURL, verifyReason);
      isSubmitted = true;
    }catch(error){
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.error || error.message;
        // console.log(message);
        toastStore.error(message);
      } else if (error instanceof Error) {
        toastStore.error(error.message);
      }
    }finally{
      isSubmitting = false;
    }
  }

  onMount(async () => {
    await checkAuth(navigate);
    const usr = await getAuthUser();
    if (usr) user = usr;

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
      const prefersDark = window.matchMedia(
        "(prefers-color-scheme: dark)"
      ).matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute(
        "data-theme",
        prefersDark ? "dark" : "light"
      );
    }
  });
</script>

<main class="premium-page">
  <div class="top-controls">
    <button
      type="button"
      class="back-button"
      aria-label="Back to home"
      on:click={() => navigate("/home")}
    >
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
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>
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
    {#if isSubmitted}
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
        <h2>Verification Request Submitted</h2>
        <p>
          Thank you for your premium verification request. Our admin team will
          review your submission and get back to you shortly.
        </p>
        <p class="note">
          Your ID information is securely encrypted and will only be used for
          verification purposes.
        </p>
        <button
          type="button"
          class="primary-button"
          on:click={() => navigate("/")}>Return to Home</button
        >
      </div>
    {:else}
      <div class="premium-info">
        <div class="premium-badge">
          <svg
            viewBox="0 0 24 24"
            width="32"
            height="32"
            fill="rgb(29, 155, 240)"
          >
            <path
              d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"
            ></path>
          </svg>
        </div>
        <h2>Get Verified with Premium</h2>
        <p>
          Join our premium membership and receive a blue checkmark next to your
          name for increased recognition and credibility in the community.
        </p>

        <div class="benefits">
          <div class="benefit-item">
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
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>Blue verification checkmark</span>
          </div>
          <div class="benefit-item">
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
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>Increased visibility in conversations</span>
          </div>
          <div class="benefit-item">
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
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>Priority customer support</span>
          </div>
        </div>
      </div>

      <div class="verification-form">
        <h3>Verification Information</h3>
        <p class="form-info">
          To become a premium user, please provide the following information.
          Our admin team will review your submission.
        </p>

        <div class="security-note">
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
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
          <span
            >Your ID information is securely encrypted and stored using
            industry-standard security protocols.</span
          >
        </div>

        <form on:submit|preventDefault={submitVerification}>
          <div class="form-group">
            <label for="id-number">National Identity Card Number</label>
            <input
              type="text"
              id="id-number"
              bind:value={idNumber}
              class:error={!idNumberValid}
              placeholder="Enter your ID number"
            />
            {#if !idNumberValid}
              <span class="error-message">Please enter a valid ID number</span>
            {/if}
          </div>

          <div class="form-group">
            <label for="verify-reason">Reason for Verification</label>
            <textarea
              id="verify-reason"
              bind:value={verifyReason}
              class:error={!reasonValid}
              placeholder="Please explain why you're requesting verification (minimum 10 characters)"
              rows="4"
            ></textarea>
            {#if !reasonValid}
              <span class="error-message">Please provide a reason for verification (minimum 10 characters)</span>
            {/if}
          </div>

          <div class="form-group">
            <label for="face-image">Upload a Photo of Your Face</label>
            <div class="image-upload-container">
              <input
                type="file"
                id="face-image"
                accept="image/*"
                on:change={handleImageUpload}
                class="hidden-input"
              />
              <label
                for="face-image"
                class="upload-button"
                class:error={!imageValid}
              >
                <span>{faceImage ? "Change Photo" : "Upload Photo"}</span>
              </label>

              {#if faceImagePreview}
                <div class="image-preview">
                  <img
                    src={faceImagePreview || "/placeholder.svg"}
                    alt="Face verification preview"
                  />
                  <button
                    type="button"
                    class="remove-image"
                    aria-label="remove button"
                    on:click={() => {
                      faceImage = null;
                      faceImagePreview = null;
                    }}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                  </button>
                </div>
              {/if}
            </div>
            {#if !imageValid}
              <span class="error-message"
                >Please upload a photo of your face</span
              >
            {/if}
          </div>

          <div class="form-actions">
            <button
              type="submit"
              class="primary-button"
              disabled={isSubmitting}
            >
              {#if isSubmitting}
                <span class="spinner"></span>
                Submitting...
              {:else}
                Submit Verification
              {/if}
            </button>
          </div>
        </form>
      </div>
    {/if}
  </div>
</main>

<style>
  .premium-page {
    width: 100%;
    height: 100vh;
    position: relative;
    overflow-y: auto;
    overflow-x: hidden;
  }

  .top-controls {
    position: absolute;
    top: 20px;
    left: 0;
    right: 0;
    width: 100%;
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    z-index: 10;
    box-sizing: border-box;
  }

  .back-button {
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

  .back-button:hover {
    background-color: var(--hover-background);
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
    max-width: 600px;
    margin: 0 auto;
    padding: 80px 20px 40px;
  }

  .premium-info {
    background-color: var(--card-background);
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 24px;
    border: 1px solid var(--border-color);
    text-align: center;
  }

  .premium-badge {
    display: flex;
    justify-content: center;
    margin-bottom: 16px;
  }

  .premium-info h2 {
    font-size: 24px;
    font-weight: 800;
    margin: 0 0 16px 0;
    color: var(--text-primary);
  }

  .premium-info p {
    font-size: 16px;
    line-height: 1.5;
    color: var(--text-primary);
    margin-bottom: 24px;
  }

  .benefits {
    display: flex;
    flex-direction: column;
    gap: 12px;
    text-align: left;
  }

  .benefit-item {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .benefit-item svg {
    color: rgb(29, 155, 240);
    flex-shrink: 0;
  }

  .verification-form {
    background-color: var(--card-background);
    border-radius: 16px;
    padding: 24px;
    border: 1px solid var(--border-color);
  }

  .verification-form h3 {
    font-size: 20px;
    font-weight: 700;
    margin: 0 0 16px 0;
    color: var(--text-primary);
  }

  .form-info {
    font-size: 15px;
    color: var(--text-primary);
    margin-bottom: 20px;
  }

  .security-note {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
    padding: 12px;
    border-radius: 8px;
    margin-bottom: 24px;
  }

  .security-note svg {
    color: var(--button-primary);
    flex-shrink: 0;
    margin-top: 2px;
  }

  .security-note span {
    font-size: 14px;
    color: var(--text-primary);
    line-height: 1.5;
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-group label {
    display: block;
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
  }

  .form-group input, .form-group textarea {
    width: 90%;
    background-color: var(--hover-background);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 12px;
    color: var(--text-primary);
    font-size: 15px;
    transition: border-color 0.2s;
  }

  .form-group textarea {
    resize: vertical;
    min-height: 100px;
    width: 90%;
    font-family: inherit;
  }

  .form-group input:focus, .form-group textarea:focus {
    border-color: var(--button-primary);
    outline: none;
  }

  .form-group input.error, .form-group textarea.error,
  .upload-button.error {
    border-color: #ff4d4f;
  }

  .error-message {
    display: block;
    color: #ff4d4f;
    font-size: 14px;
    margin-top: 4px;
  }

  .image-upload-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .hidden-input {
    display: none;
  }

  .upload-button {
    display: flex;
    gap: 8px;
    background-color: var(--hover-background);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 12px;
    cursor: pointer;
    transition: background-color 0.2s;
    width: fit-content;
  }

  .upload-button:hover {
    background-color: var(--hover-background-darker, rgba(255, 255, 255, 0.1));
  }

  .image-preview {
    position: relative;
    width: 150px;
    height: 150px;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .image-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .remove-image {
    position: absolute;
    top: 8px;
    right: 8px;
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }

  .primary-button {
    border-radius: 9999px;
    padding: 10px 20px;
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .primary-button {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
  }

  .primary-button:hover:not(:disabled) {
    background-color: var(--button-primary-hover);
  }

  .primary-button:disabled {
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
    margin-right: 8px;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .success-container {
    background-color: var(--card-background);
    border-radius: 16px;
    padding: 32px;
    text-align: center;
    border: 1px solid var(--border-color);
  }

  .success-icon {
    display: flex;
    justify-content: center;
    margin-bottom: 24px;
  }

  .success-icon svg {
    color: #52c41a;
  }

  .success-container h2 {
    font-size: 24px;
    font-weight: 800;
    margin: 0 0 16px 0;
    color: var(--text-primary);
  }

  .success-container p {
    font-size: 16px;
    line-height: 1.5;
    color: var(--text-primary);
    margin-bottom: 24px;
  }

  .success-container .note {
    font-size: 14px;
    color: var(--text-tertiary);
    margin-bottom: 32px;
  }

  @media (max-width: 768px) {
    .content {
      padding: 60px 16px 30px;
    }
  }
</style>