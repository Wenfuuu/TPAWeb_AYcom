<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { toastStore } from "../lib/stores/toast";
  import { uploadImage } from "../lib/utils/supabaseStorage";
  import type { User } from "../lib/model/user";
  import {
    fetchCommunityCategory,
    requestCreateCommunity,
    type CommunityCategory,
  } from "../lib/model/community";

  let user: User | null = null;
  let isDarkMode = true;
  let isSubmitting = false;
  let isSubmitted = false;

  let communityName = "";
  let description = "";
  let iconFile: File | null = null;
  let iconPreview: string | null = null;
  let bannerFile: File | null = null;
  let bannerPreview: string | null = null;
  let selectedCategoryID = "";
  let rule = "";

  let nameValid = true;
  let descriptionValid = true;
  let iconValid = true;
  let bannerValid = true;
  let categoryValid = true;
  let rulesValid = true;

  let availableCategories: CommunityCategory[] = [];

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }

  function handleIconUpload(event: Event) {
    const target = event.target as HTMLInputElement;
    const files = target.files;

    if (!files || files.length === 0) {
      iconFile = null;
      iconPreview = null;
      return;
    }

    const file = files[0];
    iconFile = file;

    const reader = new FileReader();
    reader.onload = (e: ProgressEvent<FileReader>) => {
      const result = e.target?.result;
      if (typeof result === "string") {
        iconPreview = result;
      }
    };

    reader.readAsDataURL(file);
  }

  function handleBannerUpload(event: Event) {
    const target = event.target as HTMLInputElement;
    const files = target.files;

    if (!files || files.length === 0) {
      bannerFile = null;
      bannerPreview = null;
      return;
    }

    const file = files[0];
    bannerFile = file;

    const reader = new FileReader();
    reader.onload = (e: ProgressEvent<FileReader>) => {
      const result = e.target?.result;
      if (typeof result === "string") {
        bannerPreview = result;
      }
    };

    reader.readAsDataURL(file);
  }

  function validateForm(): boolean {
    nameValid = communityName.trim().length >= 3 && communityName.length <= 50;

    descriptionValid =
      description.trim().length >= 10 && description.length <= 500;

    iconValid = iconFile !== null;

    bannerValid = bannerFile !== null;

    categoryValid = selectedCategoryID !== "";

    rulesValid = rule.trim().length >= 10;

    return (
      nameValid &&
      descriptionValid &&
      iconValid &&
      bannerValid &&
      categoryValid &&
      rulesValid
    );
  }

  async function submitCommunity() {
    if (!validateForm() || !user) {
      return;
    }

    isSubmitting = true;

    try {
      let iconURL = null;
      let bannerURL = null;

      if (iconFile) {
        iconURL = await uploadImage(iconFile);
      }

      if (bannerFile) {
        bannerURL = await uploadImage(bannerFile);
      }

      if (!iconURL || !bannerURL) {
        toastStore.error("Failed to upload images. Please try again.");
        isSubmitting = false;
        return;
      }

      // userid, name, desc, iconurl, bannerurl, categoryid, rule
      await requestCreateCommunity(
        user.id,
        communityName,
        description,
        iconURL,
        bannerURL,
        selectedCategoryID,
        rule
      );

      isSubmitted = true;
    } catch (error) {
        if (error instanceof Error) {
            toastStore.error(error.message);
        }
    } finally {
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

    availableCategories = await fetchCommunityCategory();
    // console.log(availableCategories);
  });
</script>

<main class="community-page">
  <div class="top-controls">
    <button
      type="button"
      class="back-button"
      aria-label="Back to home"
      on:click={() => navigate("/communities")}
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
        <h2>Community Requested Successfully</h2>
        <p>
          Your community "{communityName}" has been requested.
          <button
            type="button"
            class="primary-button"
            on:click={() => navigate("/communities")}>Go to Communities</button
          >
        </p>
      </div>
    {:else}
      <div class="community-form">
        <h3>Community Information</h3>
        <p class="form-info">
          Fill out the details below to create your community. All fields are
          required.
        </p>

        <form on:submit|preventDefault={submitCommunity}>
          <div class="form-group">
            <label for="community-name">Community Name</label>
            <input
              type="text"
              id="community-name"
              bind:value={communityName}
              class:error={!nameValid}
              placeholder="Enter community name (3-50 characters)"
            />
            {#if !nameValid}
              <span class="error-message"
                >Community name must be between 3 and 50 characters</span
              >
            {/if}
            <div class="character-count">
              {communityName.length}/50
            </div>
          </div>

          <div class="form-group">
            <label for="description">Description</label>
            <textarea
              id="description"
              bind:value={description}
              class:error={!descriptionValid}
              placeholder="Describe what your community is about (10-500 characters)"
              rows="4"
            ></textarea>
            {#if !descriptionValid}
              <span class="error-message"
                >Description must be between 10 and 500 characters</span
              >
            {/if}
            <div class="character-count">
              {description.length}/500
            </div>
          </div>

          <div class="form-group">
            <label for="community-icon">Community Icon</label>
            <div class="image-upload-container">
              <input
                type="file"
                id="community-icon"
                accept="image/*"
                on:change={handleIconUpload}
                class="hidden-input"
              />
              <label
                for="community-icon"
                class="upload-button"
                class:error={!iconValid}
              >
                <span>{iconFile ? "Change Icon" : "Upload Icon"}</span>
              </label>

              {#if iconPreview}
                <div class="image-preview icon-preview">
                  <img
                    src={iconPreview || "/placeholder.svg"}
                    alt="Community icon preview"
                  />
                  <button
                    type="button"
                    class="remove-image"
                    aria-label="remove icon"
                    on:click={() => {
                      iconFile = null;
                      iconPreview = null;
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
            {#if !iconValid}
              <span class="error-message">Please upload a community icon</span>
            {/if}
            <div class="upload-hint">
              Recommended: Square image, at least 400x400px
            </div>
          </div>

          <div class="form-group">
            <label for="community-banner">Community Banner</label>
            <div class="image-upload-container">
              <input
                type="file"
                id="community-banner"
                accept="image/*"
                on:change={handleBannerUpload}
                class="hidden-input"
              />
              <label
                for="community-banner"
                class="upload-button"
                class:error={!bannerValid}
              >
                <span>{bannerFile ? "Change Banner" : "Upload Banner"}</span>
              </label>

              {#if bannerPreview}
                <div class="image-preview banner-preview">
                  <img
                    src={bannerPreview || "/placeholder.svg"}
                    alt="Community banner preview"
                  />
                  <button
                    type="button"
                    class="remove-image"
                    aria-label="remove banner"
                    on:click={() => {
                      bannerFile = null;
                      bannerPreview = null;
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
            {#if !bannerValid}
              <span class="error-message">Please upload a community banner</span
              >
            {/if}
            <div class="upload-hint">Recommended: 1500x500px</div>
          </div>

          <div class="form-group">
            <label for="category">Category</label>
            <div class="select-wrapper">
              <select
                id="category"
                bind:value={selectedCategoryID}
                class:error={!categoryValid}
                class="styled-select"
              >
                <option value="" disabled>Select a category</option>
                {#each availableCategories as category (category.id)}
                  <option value={category.id}>{category.name}</option>
                {/each}
              </select>
              <div class="select-arrow">
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
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </div>
            </div>
            {#if !categoryValid}
              <span class="error-message">Please select a category</span>
            {/if}
            <div class="upload-hint">
              Select the category that best describes your community
            </div>
          </div>

          <div class="form-group">
            <label for="rules">Community Rules</label>
            <textarea
              id="rules"
              bind:value={rule}
              class:error={!rulesValid}
              placeholder="Define the rules members must follow in your community (minimum 10 characters)"
              rows="6"
            ></textarea>
            {#if !rulesValid}
              <span class="error-message"
                >Please add community rules (minimum 10 characters)</span
              >
            {/if}
            <div class="character-count">
              {rule.length} characters
            </div>
          </div>

          <div class="form-actions">
            <button
              type="button"
              class="secondary-button"
              on:click={() => navigate("/communities")}
            >
              Cancel
            </button>
            <button
              type="submit"
              class="primary-button"
              disabled={isSubmitting}
            >
              {#if isSubmitting}
                <span class="spinner"></span>
                Requesting...
              {:else}
                Request Create
              {/if}
            </button>
          </div>
        </form>
      </div>
    {/if}
  </div>
</main>

<style>
  .community-page {
    width: 100%;
    height: 100vh;
    position: relative;
    overflow-y: auto;
    overflow-x: hidden;
    background-color: var(--background);
    color: var(--text-primary);
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

  .community-form {
    background-color: var(--card-background);
    border-radius: 16px;
    padding: 24px;
    border: 1px solid var(--border-color);
    margin-bottom: 40px;
  }

  .community-form h3 {
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

  .form-group {
    margin-bottom: 24px;
  }

  .form-group label {
    display: block;
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
  }

  .form-group input,
  .form-group textarea,
  .form-group select {
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

  .form-group textarea {
    resize: vertical;
    min-height: 80px;
    font-family: inherit;
  }

  .form-group input:focus,
  .form-group textarea:focus,
  .form-group select:focus {
    border-color: var(--button-primary);
    outline: none;
  }

  .form-group input.error,
  .form-group textarea.error,
  .form-group select.error,
  .upload-button.error {
    border-color: #ff4d4f;
  }

  .select-wrapper {
    position: relative;
  }

  .select-arrow {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    pointer-events: none;
    color: var(--text-tertiary);
  }

  .form-group select {
    appearance: none;
    padding-right: 30px;
    cursor: pointer;
  }

  .styled-select {
    appearance: none;
    background-color: var(--hover-background);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 12px;
    color: var(--text-primary);
    font-size: 15px;
    transition:
      border-color 0.2s,
      background-color 0.2s;
    cursor: pointer;
    width: 100%;
  }

  .styled-select:focus {
    border-color: var(--button-primary);
    outline: none;
    background-color: var(--hover-background-darker);
  }

  .styled-select option {
    background-color: var(--card-background);
    color: var(--text-primary);
  }

  .error-message {
    display: block;
    color: #ff4d4f;
    font-size: 14px;
    margin-top: 4px;
  }

  .character-count {
    text-align: right;
    font-size: 12px;
    color: var(--text-tertiary);
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
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .icon-preview {
    width: 100px;
    height: 100px;
    border-radius: 50%;
  }

  .banner-preview {
    width: 100%;
    height: 150px;
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

  .upload-hint {
    font-size: 12px;
    color: var(--text-tertiary);
    margin-top: 4px;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 16px;
    margin-top: 32px;
  }

  .secondary-button {
    background: none;
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    padding: 10px 20px;
    border-radius: 9999px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .secondary-button:hover {
    background-color: var(--hover-background);
  }

  .primary-button {
    background-color: var(--button-primary);
    border: none;
    color: white;
    padding: 10px 20px;
    border-radius: 9999px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
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

  @media (max-width: 768px) {
    .content {
      padding: 60px 16px 30px;
    }
  }
</style>
