import { S3Client, GetObjectCommand, ListObjectsV2Command } from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";

// Initialize the S3 client
export const s3Client = new S3Client({
    region: import.meta.env.VITE_AWS_REGION || 'us-east-1',
    credentials: {
        accessKeyId: import.meta.env.VITE_AWS_ACCESS_KEY_ID,
        secretAccessKey: import.meta.env.VITE_AWS_SECRET_ACCESS_KEY,
    }
});

/**
 * Get a signed URL for a document in S3
 * @param s3Uri - S3 URI in the format s3://bucket-name/object-key
 * @param expiresIn - URL expiration time in seconds (default: 3600)
 */
export async function getSignedS3Url(s3Uri: string, expiresIn = 3600): Promise<string> {
    try {
        // Parse s3:// URI to get bucket and key
        const s3UriRegex = /s3:\/\/([^\/]+)\/(.+)/;
        const match = s3Uri.match(s3UriRegex);

        if (!match) {
            throw new Error(`Invalid S3 URI format: ${s3Uri}`);
        }

        const [, bucket, key] = match;

        const command = new GetObjectCommand({
            Bucket: bucket,
            Key: key,
        });

        // Generate a signed URL that expires in 1 hour
        return await getSignedUrl(s3Client, command, { expiresIn });
    } catch (error) {
        console.error("Error generating signed URL:", error);
        throw error;
    }
}

/**
 * List objects in an S3 bucket with optional prefix
 * @param bucket - S3 bucket name
 * @param prefix - Object key prefix (folder path)
 * @param maxKeys - Maximum number of objects to return
 */
export async function listS3Objects(bucket: string, prefix = '', maxKeys = 1000) {
    try {
        const command = new ListObjectsV2Command({
            Bucket: bucket,
            Prefix: prefix,
            MaxKeys: maxKeys,
        });

        const response = await s3Client.send(command);
        return response.Contents || [];
    } catch (error) {
        console.error("Error listing S3 objects:", error);
        throw error;
    }
}