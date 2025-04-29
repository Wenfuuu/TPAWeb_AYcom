<script lang="ts">
//   import { onMount } from "svelte";

  // Sample posts data
  const posts = [
    {
      id: 1,
      user: {
        name: "Desdee",
        username: "@NormalSmith",
        verified: true,
        avatar: "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/placeholder-ob7miW3mUreePYfXdVwkpFWHthzoR5.svg?height=48&width=48"
      },
      date: "Apr 1",
      content: "Hai sederhana apa yang bikin kalian emosi?",
      images: [
        "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/placeholder-ob7miW3mUreePYfXdVwkpFWHthzoR5.svg?height=200&width=250",
        "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/placeholder-ob7miW3mUreePYfXdVwkpFWHthzoR5.svg?height=200&width=250",
        "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/placeholder-ob7miW3mUreePYfXdVwkpFWHthzoR5.svg?height=200&width=250",
        "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/placeholder-ob7miW3mUreePYfXdVwkpFWHthzoR5.svg?height=200&width=250"
      ],
      caption: "sederhana tapi aku marah",
      stats: {
        replies: 221,
        reposts: 33,
        likes: 2.8,
        views: 450
      }
    },
    {
      id: 2,
      user: {
        name: "BAZNAS",
        username: "@baznasIndonesia",
        verified: true,
        avatar: "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/placeholder-ob7miW3mUreePYfXdVwkpFWHthzoR5.svg?height=48&width=48"
      },
      date: "2h",
      content: "Mari tunaikan kewajiban zakat untuk membantu masyarakat Indonesia",
      images: [],
      stats: {
        replies: 5,
        reposts: 12,
        likes: 45,
        views: 320
      }
    }
  ];

  let newPostText = "";

  function handlePost() {
    if (newPostText.trim()) {
      console.log("Posting:", newPostText);
      newPostText = "";
    }
  }
</script>

<main class="content-area">
  <div class="header">
    <div class="tabs">
      <button type="button" class="tab active">For you</button>
      <button type="button" class="tab">Following</button>
    </div>
  </div>

  <div class="compose-post">
    <div class="compose-avatar">
      <img src="https://hebbkx1anhila5yf.public.blob.vercel-storage.com/placeholder-ob7miW3mUreePYfXdVwkpFWHthzoR5.svg?height=48&width=48" alt="Your avatar" />
    </div>
    <div class="compose-input">
      <textarea
        placeholder="What's happening?"
        bind:value={newPostText}
      ></textarea>
      <div class="compose-actions">
        <div class="compose-tools">
          <button type="button" class="tool-button" aria-label="Tool Button">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect><circle cx="9" cy="9" r="2"></circle><path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path></svg>
          </button>
          <button type="button" class="tool-button" aria-label="Tool Button">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"></path><polyline points="14 2 14 8 20 8"></polyline></svg>
          </button>
          <button type="button" class="tool-button" aria-label="Tool Button">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 8c-2.8 0-5 2.2-5 5s2.2 5 5 5 5-2.2 5-5-2.2-5-5-5Z"></path><path d="M12 3a2 2 0 0 0 2-2H10a2 2 0 0 0 2 2Z"></path><path d="M12 19v2"></path><path d="m5 8-3-3"></path><path d="M19 8h3"></path><path d="m16.95 19.5 2.5 2.5"></path><path d="m7.05 19.5-2.5 2.5"></path></svg>
          </button>
          <button type="button" class="tool-button" aria-label="Tool Button">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M8 14s1.5 2 4 2 4-2 4-2"></path><line x1="9" x2="9.01" y1="9" y2="9"></line><line x1="15" x2="15.01" y1="9" y2="9"></line></svg>
          </button>
          <button type="button" class="tool-button" aria-label="Tool Button">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="2"></rect><path d="M8 12h8"></path><path d="M12 8v8"></path></svg>
          </button>
        </div>
        <button
        type="button"
        class="post-button"
        disabled={!newPostText.trim()}
        on:click={handlePost}
        >
          Post
        </button>
      </div>
    </div>
  </div>

  <div class="posts">
    {#each posts as post (post.id)}
      <div class="post">
        <div class="post-avatar">
          <img src={post.user.avatar || "/placeholder.svg"} alt={post.user.name} />
        </div>
        <div class="post-content">
          <div class="post-header">
            <div class="post-user-info">
              <span class="post-user-name">{post.user.name}</span>
              {#if post.user.verified}
                <span class="verified-badge">
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="rgb(29, 155, 240)">
                    <path d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"></path>
                  </svg>
                </span>
              {/if}
              <span class="post-username">{post.user.username}</span>
              <span class="post-date">· {post.date}</span>
            </div>
            <button type="button" class="post-more" aria-label="Post More">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="1"></circle>
                <circle cx="19" cy="12" r="1"></circle>
                <circle cx="5" cy="12" r="1"></circle>
              </svg>
            </button>
          </div>

          <div class="post-text">{post.content}</div>

          {#if post.images && post.images.length > 0}
            <div class="post-images" class:grid-2={post.images.length === 2} class:grid-3={post.images.length === 3} class:grid-4={post.images.length === 4}>
              {#each post.images as image, i (i)}
                <div class="image-container">
                  <img src={image || "/placeholder.svg"} alt="Post image {i+1}" />
                </div>
              {/each}
            </div>
            {#if post.caption}
              <div class="post-caption">{post.caption}</div>
            {/if}
          {/if}

          <div class="post-actions">
            <button type="button" class="post-action">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
              </svg>
              <span>{post.stats.replies}</span>
            </button>
            <button type="button" class="post-action">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m17 2 4 4-4 4"></path>
                <path d="M3 11v-1a4 4 0 0 1 4-4h14"></path>
                <path d="m7 22-4-4 4-4"></path>
                <path d="M21 13v1a4 4 0 0 1-4 4H3"></path>
              </svg>
              <span>{post.stats.reposts}</span>
            </button>
            <button type="button" class="post-action">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"></path>
              </svg>
              <span>{post.stats.likes}K</span>
            </button>
            <button type="button" class="post-action" aria-label="View count">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z"></path>
                <path d="M12 9a3 3 0 1 0 0 6 3 3 0 1 0 0-6z"></path>
              </svg>
              <span>{post.stats.views}K</span>
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

<style>
  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    /* min-width: 50%; */
    /* max-width: 600px; */
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

  .compose-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
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