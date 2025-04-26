import { createClient } from "@supabase/supabase-js";

// Initialize Supabase client
const supabaseUrl = "https://nyijguoyjemxmdyzzdry.supabase.co";
const supabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im55aWpndW95amVteG1keXp6ZHJ5Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDU2NTI4MTMsImV4cCI6MjA2MTIyODgxM30.cGZzeD6ojLZjxPu9SpFGkjnkPCAnUNQy8GFCXkLD8uk";
const supabase = createClient(supabaseUrl, supabaseKey);

/**
 * Upload a file to Supabase Storage during registration
 * @param file The file to upload
 * @param fileType The type of file ('profile' or 'banner')
 * @returns The URL of the uploaded file or null if upload failed
 */
export async function uploadRegistrationFile(
  file: File,
  fileType: "profile" | "banner"
): Promise<string | null> {
  try {
    // Generate a unique filename
    const fileExt = file.name.split(".").pop();
    const timestamp = new Date().getTime();
    const randomString = Math.random().toString(36).substring(2, 10);
    const fileName = `${fileType}-${timestamp}-${randomString}.${fileExt}`;

    // Upload the file to the registration-uploads bucket
    const { data, error } = await supabase.storage
      .from("registration-uploads")
      .upload(fileName, file, {
        cacheControl: "3600",
        upsert: false
      });

    if (error) {
      console.error("Error uploading file:", error);
      throw error;
    }

    // Get the public URL
    const { data: urlData } = supabase.storage
      .from("registration-uploads")
      .getPublicUrl(fileName);

    return urlData.publicUrl;
  } catch (error) {
    console.error("Error in uploadRegistrationFile:", error);
    return null;
  }
}

/**
 * Upload a profile picture during registration
 * @param file The image file
 * @returns The URL of the uploaded profile picture
 */
export async function uploadProfilePicture(file: File): Promise<string | null> {
  return uploadRegistrationFile(file, "profile");
}

/**
 * Upload a banner image during registration
 * @param file The image file
 * @returns The URL of the uploaded banner
 */
export async function uploadBannerImage(file: File): Promise<string | null> {
  return uploadRegistrationFile(file, "banner");
}