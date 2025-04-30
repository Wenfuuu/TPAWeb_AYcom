<script lang="ts">
  import { onMount } from "svelte";
  import { type User, type UserThread } from "../model/user";
  import { uploadImage } from "../utils/supabaseStorage";
  import axios from "axios";
  import { fetchCommunity, fetchThreadCategory, fetchThreadWithUser, type Community, type Thread, type ThreadCategory } from "../model/thread";

  // let threads: Thread[] = [];
  let threads: (Thread & { userData: UserThread })[] = [];
  // let users: User[] = [];
  let threadCategories: ThreadCategory[] = [];
  let communities: Community[] = [];
  // let threadFiles: ThreadFile[] = [];

  export let user: User | null;
  let newPostText = "";
  let uploadedImages: string[] = [];
  let imageFiles: File[] = [];
  let supabaseURLS: string[] = [];
  let currentImageIndex = 0;

  // New variables for category and community selection
  let selectedCategoryId = "";
  let isPostingToCommunity = false;
  let selectedCommunityId = "";
  let showCategoryDropdown = false;
  let showCommunityDropdown = false;

  type ReplySettingType = "everyone" | "following" | "verified";
  let replySettings: ReplySettingType = "everyone";
  let showReplyOptions = false;

  onMount(async() => {
    // const user = await getAuthUser();
    // if(user) {
    //   currentUser = user;
    // }

    threads = await fetchThreadWithUser();
    // users = await fetchUser();
    threadCategories = await fetchThreadCategory();
    communities = await fetchCommunity();
    // threadFiles = await fetchFiles();

    console.log(threads);
    // console.log(users);
    console.log(threadCategories);
    console.log(communities);
    // console.log(threadFiles);
  });

  async function handlePost(): Promise<void> {
    if (newPostText.trim()) {
      if (!selectedCategoryId) {
        alert("Please select a category for your thread");
        return;
      }

      if (isPostingToCommunity && !selectedCommunityId) {
        alert("Please select a community to post to");
        return;
      }

      console.log("Posting:", newPostText, uploadedImages);

      // upload images to supabase
      if (imageFiles.length > 0) {
        await uploadImagesToSupabase();
      }

      // call api to create thread
      let isAdvertisement;
      if(user?.role == "user") isAdvertisement = false;
      else isAdvertisement = true;
      // console.log(currentUser.id);
      await axios.post("http://localhost:8080/api/threads/post", {
        userId: user?.id,
        categoryId: selectedCategoryId,
        communityId: selectedCommunityId,
        content: newPostText,
        replyPermission: replySettings,
        isAdvertisement,
        scheduledAt: "2025-04-30T15:04:00Z",
        imageUrls: supabaseURLS
      });

      newPostText = "";
      uploadedImages = [];
      imageFiles = [];
      supabaseURLS = [];
      selectedCategoryId = "";
      selectedCommunityId = "";
      isPostingToCommunity = false;
    }
  }

  // upload and get urls
  async function uploadImagesToSupabase(): Promise<boolean> {
    if (!imageFiles || imageFiles.length === 0) return false;

    // upload all image here
    for (const file of imageFiles) {
      const supaURL = await uploadImage(file);
      if (supaURL) {
        supabaseURLS.push(supaURL);
      }
    }

    return true;
  }

  // for image preview
  function handleImageUpload(event: Event): void {
    const target = event.target as HTMLInputElement;
    const files = target.files;

    if (!files || files.length === 0) return;

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      imageFiles = [...imageFiles, file];
      const reader = new FileReader();

      reader.onload = (e: ProgressEvent<FileReader>) => {
        const result = e.target?.result;
        if (typeof result === "string") {
          uploadedImages = [...uploadedImages, result];
        }
      };

      reader.readAsDataURL(file);
    }
  }

  function removeImage(index: number): void {
    uploadedImages = uploadedImages.filter((_, i) => i !== index);
    if (currentImageIndex >= uploadedImages.length) {
      currentImageIndex = Math.max(0, uploadedImages.length - 1);
    }
  }

  function toggleReplyOptions(): void {
    showReplyOptions = !showReplyOptions;
  }

  function setReplyOption(option: ReplySettingType): void {
    replySettings = option;
    showReplyOptions = false;
  }

  // New functions for category and community selection
  function toggleCategoryDropdown(): void {
    showCategoryDropdown = !showCategoryDropdown;
    // Close other dropdowns
    showReplyOptions = false;
    showCommunityDropdown = false;
  }

  function selectCategory(categoryId: string): void {
    selectedCategoryId = categoryId;
    showCategoryDropdown = false;
  }

  function toggleCommunityDropdown(): void {
    showCommunityDropdown = !showCommunityDropdown;
    // Close other dropdowns
    showReplyOptions = false;
    showCategoryDropdown = false;
  }

  function selectCommunity(communityId: string): void {
    selectedCommunityId = communityId;
    showCommunityDropdown = false;
  }

  function togglePostingType(): void {
    isPostingToCommunity = !isPostingToCommunity;
    if (!isPostingToCommunity) {
      selectedCommunityId = "";
    }
  }
</script>

{#if user}
  <main class="content-area">
  <div class="header">
    <div class="tabs">
      <button type="button" class="tab active">For you</button>
      <button type="button" class="tab">Following</button>
    </div>
  </div>

  <div class="compose-post">
    <div class="compose-avatar">
      <img src={user.profile_picture} alt="Your avatar" />
    </div>
    <div class="compose-input">
      <textarea
        placeholder="What's happening?"
        bind:value={newPostText}
      ></textarea>

      {#if uploadedImages.length > 0}
        <div class="image-carousel">
          <div class="carousel-container">
            <div class="carousel-images" class:has-multiple={uploadedImages.length > 1}>
              {#each uploadedImages as image, i (i)}
                <div class="carousel-image" class:active={i === currentImageIndex}>
                  <img src={image || "/placeholder.svg"} alt="Uploaded image {i+1}" />
                  <button type="button" class="edit-button">Edit</button>
                  <button type="button" class="remove-button" aria-label="Remove Button" on:click={() => removeImage(i)}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                  </button>
                </div>
              {/each}
            </div>

            {#if uploadedImages.length > 1}
              <div class="carousel-nav">
                <button type="button" class="nav-button prev" aria-label="Button Prev" disabled={currentImageIndex === 0} on:click={() => currentImageIndex = Math.max(0, currentImageIndex - 1)}>
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="m15 18-6-6 6-6"></path>
                  </svg>
                </button>
                <button type="button" class="nav-button next" aria-label="Button Next" disabled={currentImageIndex === uploadedImages.length - 1} on:click={() => currentImageIndex = Math.min(uploadedImages.length - 1, currentImageIndex + 1)}>
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="m9 18 6-6-6-6"></path>
                  </svg>
                </button>
              </div>

              <div class="carousel-indicators">
                {#each uploadedImages as _, i (i)}
                  <span class="indicator" class:active={i === currentImageIndex}></span>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      {/if}

      <!-- new stuff -->
       <!-- Post options section -->
      <div class="post-options">
        <!-- Category selection -->
        <div class="option-dropdown">
          <button type="button" class="option-button" on:click={toggleCategoryDropdown}>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z"></path>
            </svg>
            {selectedCategoryId ? threadCategories.find(c => c.id === selectedCategoryId)?.name || "Select category" : "Select category"}
          </button>

          {#if showCategoryDropdown}
            <div class="dropdown-menu">
              {#each threadCategories as category (category.id)}
                <button type="button" class="dropdown-item" on:click={() => selectCategory(category.id)}>
                  <div class="item-text">
                    <h4>{category.name}</h4>
                  </div>
                  {#if selectedCategoryId === category.id}
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>

        <!-- Post type toggle -->
        <div class="post-type-toggle">
          <button type="button" class="option-button" on:click={togglePostingType}>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            {isPostingToCommunity ? "Posting to community" : "Posting personally"}
          </button>
        </div>

        <!-- Community selection (only shown when posting to community) -->
        {#if isPostingToCommunity}
          <div class="option-dropdown">
            <button type="button" class="option-button" on:click={toggleCommunityDropdown}>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5Z"></path>
                <path d="M3 10h18"></path>
                <path d="M10 3v18"></path>
              </svg>
              {selectedCommunityId ? communities.find(c => c.id === selectedCommunityId)?.name || "Select community" : "Select community"}
            </button>

            {#if showCommunityDropdown}
              <div class="dropdown-menu">
                {#each communities as community (community.id)}
                  <button type="button" class="dropdown-item" on:click={() => selectCommunity(community.id)}>
                    <div class="item-text">
                      <h4>{community.name}</h4>
                      <p>{community.description}</p>
                    </div>
                    {#if selectedCommunityId === community.id}
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    {/if}
                  </button>
                {/each}
              </div>
            {/if}
          </div>
        {/if}
      </div>

      <div class="reply-settings">
        <button type="button" class="reply-button" on:click={toggleReplyOptions}>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="M12 8h.01"></path>
            <path d="M11 12h1v4h1"></path>
          </svg>
          {#if replySettings === "everyone"}
            Everyone can reply
          {:else if replySettings === "following"}
            People you follow can reply
          {:else}
            Only people you mention can reply
          {/if}
        </button>

        {#if showReplyOptions}
          <div class="reply-dropdown">
            <button type="button" class="reply-option" on:click={() => setReplyOption("everyone")}>
              <div class="option-icon everyone" class:selected={replySettings === "everyone"}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                </svg>
              </div>
              <div class="option-text">
                <h4>Everyone</h4>
                <p>Anyone can reply to this post</p>
              </div>
            </button>

            <button type="button" class="reply-option" on:click={() => setReplyOption("following")}>
              <div class="option-icon following" class:selected={replySettings === "following"}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                </svg>
              </div>
              <div class="option-text">
                <h4>People you follow</h4>
                <p>Only people you follow can reply</p>
              </div>
            </button>

            <button type="button" class="reply-option" on:click={() => setReplyOption("verified")}>
              <div class="option-icon mentioned" class:selected={replySettings === "verified"}>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>
                </svg>
              </div>
              <div class="option-text">
                <h4>Verified only</h4>
                <p>Only verified accounts can reply</p>
              </div>
            </button>
          </div>
        {/if}
      </div>

      <div class="compose-actions">
        <div class="compose-tools">
          <!-- Image upload button -->
          <label class="tool-button" for="image-upload">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
              <circle cx="9" cy="9" r="2"></circle>
              <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
            </svg>
            <input
              type="file"
              id="image-upload"
              accept="image/*"
              multiple
              on:change={handleImageUpload}
              class="hidden-input"
            />
          </label>

          <!-- Poll button -->
          <button type="button" class="tool-button" aria-label="Create poll">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 3v18h18"></path>
              <path d="M18 17V9"></path>
              <path d="M13 17V5"></path>
              <path d="M8 17v-3"></path>
            </svg>
          </button>

          <!-- Schedule button -->
          <button type="button" class="tool-button" aria-label="Schedule post">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </button>
        </div>

        <button
          type="button"
          class="post-button"
          disabled={!newPostText.trim() && uploadedImages.length === 0}
          on:click={handlePost}
        >
          Post
        </button>
      </div>
    </div>
  </div>

  <div class="posts">
    {#each threads as thread (thread.id)}
      <div class="post">
        <div class="post-avatar">
          <img src={thread.userData.profile || "/placeholder.svg"} alt={thread.userData.name} />
        </div>
        <div class="post-content">
          <div class="post-header">
            <div class="post-user-info">
              <span class="post-user-name">{thread.userData.name}</span>
              {#if thread.userData.isPremium}
                <span class="verified-badge">
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="rgb(29, 155, 240)">
                    <path d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"></path>
                  </svg>
                </span>
              {/if}
              <span class="post-username">{thread.userData.username}</span>
              <span class="post-date">· {thread.created_at}</span>
            </div>
            <button type="button" class="post-more" aria-label="Post More">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="1"></circle>
                <circle cx="19" cy="12" r="1"></circle>
                <circle cx="5" cy="12" r="1"></circle>
              </svg>
            </button>
          </div>

          {#if thread.community_id}
            <div class="post-community">
              Posted in {communities.find(c => c.id === thread.community_id)?.name || "Community"}
            </div>
          {/if}

          <div class="post-text">{thread.content}</div>

          {#if thread.images && thread.images.length > 0}
            <div class="post-images" class:grid-2={thread.images.length === 2} class:grid-3={thread.images.length === 3} class:grid-4={thread.images.length === 4}>
              {#each thread.images as image, i (i)}
                <div class="image-container">
                  <img src={image || "/placeholder.svg"} alt="Post image {i+1}" />
                </div>
              {/each}
            </div>
            {#if thread.content}
              <div class="post-caption">{thread.content}</div>
            {/if}
          {/if}

          <div class="post-actions">
            <button type="button" class="post-action">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
              </svg>
              <span>{thread.insights.replies}</span>
            </button>
            <button type="button" class="post-action">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m17 2 4 4-4 4"></path>
                <path d="M3 11v-1a4 4 0 0 1 4-4h14"></path>
                <path d="m7 22-4-4 4-4"></path>
                <path d="M21 13v1a4 4 0 0 1-4 4H3"></path>
              </svg>
              <span>{thread.insights.reposts}</span>
            </button>
            <button type="button" class="post-action">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"></path>
              </svg>
              <span>{thread.insights.likes}K</span>
            </button>
            <button type="button" class="post-action" aria-label="Post Action">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="17 8 12 3 7 8"></polyline>
                <line x1="12" x2="12" y1="3" y2="15"></line>
              </svg>
            </button>
          </div>
        </div>
      </div>
    {/each}
  </div>
</main>
{/if}

<style>
  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 50%;
    height: 100vh;
    overflow-y: auto;
  }

  .header {
    position: sticky;
    top: 0;
    background-color: var(--card-background);
    backdrop-filter: blur(12px);
    z-index: 10;
    border-bottom: 1px solid var(--border-color);
  }

  .tabs {
    display: flex;
    height: 53px;
  }

  .tab {
    flex: 1;
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 15px;
    font-weight: 500;
    padding: 16px;
    cursor: pointer;
    position: relative;
    transition: background-color 0.2s;
  }

  .tab:hover {
    background-color: var(--hover-background);
  }

  .tab.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 56px;
    height: 4px;
    background-color: var(--button-primary);
    border-radius: 9999px;
  }

  .compose-post {
    display: flex;
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .compose-avatar {
    margin-right: 12px;
  }

  .compose-avatar img {
    width: 48px;
    height: 48px;
    border-radius: 50%;
  }

  .compose-input {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  textarea {
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 20px;
    font-family: inherit;
    height: 52px;
    padding: 12px 0;
    resize: none;
    outline: none;
  }

  textarea::placeholder {
    color: var(--text-tertiary);
  }

  /* Hide file input but keep it accessible */
  .hidden-input {
    display: none;
  }

  /* Post options */
  .post-options {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin: 12px 0;
  }

  .option-dropdown {
    position: relative;
  }

  .option-button {
    display: flex;
    align-items: center;
    gap: 8px;
    background: none;
    border: 1px solid var(--border-color);
    border-radius: 16px;
    color: var(--button-primary);
    font-size: 14px;
    padding: 6px 12px;
    cursor: pointer;
  }

  .dropdown-menu {
    position: absolute;
    top: 100%;
    left: 0;
    width: 100%;
    background-color: var(--card-background);
    border-radius: 16px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    z-index: 100;
    margin-top: 8px;
    padding: 12px;
    border: 1px solid var(--border-color);
  }

  .dropdown-item {
    width: 100%;
    background-color: var(--background);
    border: none;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 12px;
    cursor: pointer;
    border-radius: 8px;
    transition: background-color 0.2s;
    text-align: left;
  }

  .dropdown-item:hover {
    background-color: var(--hover-background);
  }

  .item-text h4 {
    margin: 0;
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .item-text p {
    margin: 4px 0 0 0;
    font-size: 13px;
    color: var(--text-tertiary);
  }

  .post-community {
    color: var(--button-primary);
    font-size: 14px;
    margin-bottom: 8px;
  }

  /* Image carousel styles */
  .image-carousel {
    margin: 12px 0;
    border-radius: 16px;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .carousel-container {
    position: relative;
    width: 100%;
    height: 280px;
    background-color: #000;
    overflow: hidden;
  }

  .carousel-images {
    display: flex;
    width: 100%;
    height: 100%;
  }

  .carousel-images.has-multiple {
    flex-wrap: nowrap;
  }

  .carousel-image {
    position: relative;
    min-width: 100%;
    height: 100%;
    display: none;
  }

  .carousel-image.active {
    display: block;
  }

  .carousel-image img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .edit-button {
    position: absolute;
    top: 8px;
    right: 8px;
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 9999px;
    padding: 6px 12px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
  }

  .remove-button {
    position: absolute;
    top: 8px;
    left: 8px;
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 50%;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .carousel-nav {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
  }

  .nav-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    pointer-events: auto;
  }

  .nav-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .nav-button.prev {
    left: 8px;
  }

  .nav-button.next {
    right: 8px;
  }

  .carousel-indicators {
    position: absolute;
    bottom: 12px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 4px;
  }

  .indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.5);
  }

  .indicator.active {
    background-color: white;
  }

  /* Reply settings styles */
  .reply-settings {
    position: relative;
    margin: 12px 0;
  }

  .reply-button {
    display: flex;
    align-items: center;
    gap: 8px;
    background: none;
    border: none;
    color: var(--button-primary);
    font-size: 14px;
    cursor: pointer;
    padding: 4px 0;
  }

  .reply-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    width: 300px;
    background-color: var(--card-background);
    border-radius: 16px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    z-index: 100;
    margin-top: 8px;
    padding: 12px;
    border: 1px solid var(--border-color);
  }

  .reply-option {
    width: 100%;
    background-color: var(--background);
    border: none;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px;
    cursor: pointer;
    border-radius: 8px;
    transition: background-color 0.2s;
    text-align: left;
  }

  .reply-option:hover {
    background-color: var(--hover-background);
  }

  .option-icon {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: var(--button-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
  }

  .option-icon.selected {
    background-color: var(--button-primary);
  }

  .option-text h4 {
    margin: 0;
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .option-text p {
    margin: 4px 0 0 0;
    font-size: 13px;
    color: var(--text-tertiary);
  }

  .compose-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
    border-top: 1px solid var(--border-color);
    padding-top: 12px;
  }

  .compose-tools {
    display: flex;
    gap: 4px;
  }

  .tool-button {
    background: none;
    border: none;
    color: var(--button-primary);
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .tool-button:hover {
    background-color: var(--hover-background);
  }

  .post-button {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .post-button:hover:not(:disabled) {
    background-color: var(--button-primary-hover);
  }

  .post-button:disabled {
    background-color: var(--button-primary-hover);
  }

  .post-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .posts {
    display: flex;
    flex-direction: column;
  }

  .post {
    display: flex;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color);
    transition: background-color 0.2s;
  }

  .post:hover {
    background-color: var(--hover-background);
  }

  .post-avatar {
    margin-right: 12px;
  }

  .post-avatar img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
  }

  .post-content {
    flex: 1;
    min-width: 0;
  }

  .post-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 4px;
  }

  .post-user-info {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    min-width: 0;
  }

  .post-user-name {
    font-weight: 700;
    color: var(--text-primary);
    margin-right: 4px;
    white-space: nowrap;
  }

  .verified-badge {
    margin-right: 4px;
    display: flex;
    align-items: center;
  }

  .post-username, .post-date {
    color: var(--text-tertiary);
    white-space: nowrap;
  }

  .post-more {
    background: none;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s, color 0.2s;
  }

  .post-more:hover {
    background-color: rgba(29, 155, 240, 0.1);
    color: var(--button-primary);
  }

  .post-text {
    color: var(--text-primary);
    margin-bottom: 12px;
    word-wrap: break-word;
  }

  .post-images {
    border-radius: 16px;
    overflow: hidden;
    margin-bottom: 12px;
    display: grid;
    gap: 2px;
    max-height: 510px;
  }

  .post-images.grid-2 {
    grid-template-columns: repeat(2, 1fr);
  }

  .post-images.grid-3 {
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr 1fr;
  }

  .post-images.grid-3 .image-container:first-child {
    grid-row: span 2;
  }

  .post-images.grid-4 {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
  }

  .image-container {
    height: 100%;
    min-height: 100px;
    overflow: hidden;
  }

  .image-container img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .post-caption {
    color: var(--text-primary);
    margin-bottom: 12px;
  }

  .post-actions {
    display: flex;
    justify-content: space-between;
    max-width: 425px;
  }

  .post-action {
    background: none;
    border: none;
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    transition: color 0.2s;
  }

  .post-action:hover {
    color: var(--button-primary);
  }

  .post-action span {
    font-size: 13px;
  }

  @media (max-width: 1200px) {
    .content-area {
      width: 60%;
    }
  }

  @media (max-width: 768px) {
    .content-area {
      width: 100%;
      margin-bottom: 60px;
    }
  }
</style>