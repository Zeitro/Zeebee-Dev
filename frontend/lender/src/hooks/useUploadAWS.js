import React, { useState } from 'react';
import AWS from 'aws-sdk';
// AWS.config.update({
// 	region:  process.env.NEXT_PUBLIC_AWS_REGION,
// 	credentials: new AWS.Credentials({
// 		accessKeyId: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY_ID,
// 		secretAccessKey: process.env.NEXT_PUBLIC_AWS_SECRET_ACCESS_KEY,
// 	}),
// });
// 文件上传亚马逊服务器
const useUploadAWS = () => {
	const [progress, setProgress] = useState(0);
	const [error, setError] = useState(null);
	const [uploading, setUploading] = useState(false);
	const s3 = new AWS.S3({
		region: process.env.NEXT_PUBLIC_AWS_REGION,
		accessKeyId: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY_ID,
		secretAccessKey: process.env.NEXT_PUBLIC_AWS_SECRET_ACCESS_KEY,
	});

	const uploadFile = async (file) => {
		const params = {
			Bucket: process.env.NEXT_PUBLIC_BUCKET_NAME,
			Key: `${process.env.NEXT_PUBLIC_AWS_FILE_FOLDER}/pdf/${Date.now()}_${file.name}`,
			Body: file,
			// ContentType: file.type,
		};

		try {
			setUploading(true);
			const upload = s3.upload(params);
			upload.on('httpUploadProgress', (progress) => {
				setProgress((progress.loaded / progress.total) * 100);
			});
			await upload.promise();
			setUploading(false);
		} catch (error) {
			setError(error);
			setUploading(false);
		}
	};

	return {
		uploadFile,
		progress,
		error,
		uploading,
	};
};
export default useUploadAWS;
