Conversation with Cloudinary
Started on October 30, 2025 at 10:52 AM Central Time (US & Canada) time CDT (GMT-0500)

---

10:52 AM | Ask Cloudinary AI: Hi Olivier! This is a bot speaking. I’m here to answer your questions on Cloudinary.

I'm trained on Cloudinary Documentation (https://cloudinary.com/documentation) and still learning so make sure to let me know if I helped.

10:52 AM | Ask Cloudinary AI: How can I help?

10:52 AM | Olivier Coq: public id and signature for uploading via POST?

10:52 AM | Ask Cloudinary AI: For POST uploads, the public_id is optional - if not specified, we'll auto-generate one for you. (https://cloudinary.com/documentation/image_upload_api_reference) 
 The signature parameter is required for signed REST API calls and authenticates your request based on the parameters you use. 
 To manually generate a signature: 
 1. Create a string with all POST parameters (except file, cloud_name, resource_type, and api_key) 
2. Add the timestamp parameter 
3. Sort parameters alphabetically and format as param=value&param=value 
4. Append your API secret 
5. Create an SHA hash of the complete string (https://cloudinary.com/documentation/authentication_signatures) 
 Example: For API secret abcd, API key 1234, timestamp 1315060510, and uploading from a URL, the signature would be bfd09f95f331f558cbd1320e67aa8d488770583e. 
 If you're using our SDKs, signatures are automatically generated for you! 
 
Sources:
Upload API Reference | Documentation (https://cloudinary.com/documentation/image_upload_api_reference)
Generating authentication signatures | Documentation (https://cloudinary.com/documentation/authentication_signatures)

10:53 AM | Ask Cloudinary AI: Is that what you were looking for?

10:55 AM | Olivier Coq: tried signed upload. got error: "Invalid Signature 61ff76640a4a23242f08e4e45d53b9d144222e06. String to sign - 'timestamp=1761839719'."

10:55 AM | Ask Cloudinary AI: The "Invalid Signature" error with only timestamp=1761839719 in your string to sign indicates you're missing required parameters for signature generation. (https://cloudinary.com/documentation/signatures) 
 For signed uploads, your signature string must include: 
 - All POST parameters except file, cloud_name, resource_type, and api_key 
- The timestamp parameter 
- Parameters sorted alphabetically in param=value&param=value format 
 Your current string only has the timestamp, but you need to include other parameters like public_id, transformations, or any upload options you're using. 
 Another common cause is SHA algorithm mismatch. If your account is configured for SHA-256 only, ensure you're using the correct algorithm for signature generation. (https://support.cloudinary.com/hc/en-us/articles/10679615558802-Why-am-I-getting-Invalid-Signature-errors-when-using-Cloudinary-SDK-Upload-API-methods) 
 To better troubleshoot, could you share: 
 - Your complete parameter list for the upload 
- How you're generating the signature 
- Which SDK/method you're using 
 
Sources:
Why am I getting "Invalid Signature" errors when using Cloudinary SDK Upload API methods? – Cloudinary Support (https://support.cloudinary.com/hc/en-us/articles/10679615558802-Why-am-I-getting-Invalid-Signature-errors-when-using-Cloudinary-SDK-Upload-API-methods)
Generating authentication signatures | Documentation (https://cloudinary.com/documentation/authentication_signatures)
Cloudinary Signatures | Documentation (https://cloudinary.com/documentation/signatures)

10:56 AM | Olivier Coq: where to find public_id

10:56 AM | Ask Cloudinary AI: You can find the public_id in several ways: 
 After uploading: The public_id is returned in the upload API response when you upload an asset. (https://cloudinary.com/documentation/upload_parameters) 
 In Media Library: Browse to Media Library (https://console.cloudinary.com/console/media_library/search), double-click any asset to open the Manage page, navigate to the Summary tab, and copy the Public ID. (https://cloudinary.com/documentation/javascript_quick_start) 
 Setting your own: You can specify a custom public_id during upload, or use use_filename: true to use the original filename as the public_id. (https://cloudinary.com/documentation/image_upload_api_reference) 
 Auto-generated: If you don't specify one, we'll generate a random public_id like 8jsb1xofxdqamu2rzwt9q. 
 For your signature error, you'll need the public_id (if you're setting one) to include in your signature string along with other parameters. 
 
Sources:
Customizing uploads | Documentation (https://cloudinary.com/documentation/upload_parameters)
Upload API Reference | Documentation (https://cloudinary.com/documentation/image_upload_api_reference)
JavaScript quick start | Documentation (https://cloudinary.com/documentation/javascript_quick_start)

10:56 AM | Ask Cloudinary AI: Was that helpful?

10:57 AM | Olivier Coq: Error: "Invalid Signature f077750b38c73a6e2af8adc799e4610d194eb2c9. String to sign - 'public_id=8jsb1xofxdqamu2rzwt9q&timestamp=1761839827'."

10:57 AM | Ask Cloudinary AI: Thanks for sharing the error details! Could you please tell me how you are generating the signature and which SDK or method you are using? This will help me understand the issue better and assist you further.

10:57 AM | Olivier Coq: testing signed upload feature in Postman

10:58 AM | Ask Cloudinary AI: Thanks for the info! Could you share the exact parameters you're including in Postman and how you're generating the signature? That way, I can help you troubleshoot this signed upload issue more effectively.

11:00 AM | Olivier Coq: example string for manually creating signatures, please

11:01 AM | Ask Cloudinary AI: Here's a complete example for manually creating signatures: 
 Parameters to sign: 
 - timestamp: 1315060510 
- public_id: sample_image 
- eager: w_400,h_300,c_pad|w_260,h_200,c_crop 
 Step-by-step process: 
 1. Create parameter string (alphabetically sorted):
eager=w_400,h_300,c_pad|w_260,h_200,c_crop&public_id=sample_image&timestamp=1315060510 
2. Append API secret (example: abcd):
eager=w_400,h_300,c_pad|w_260,h_200,c_crop&public_id=sample_image&timestamp=1315060510abcd 
3. Generate SHA-1 hash:
Result: bfd09f95f331f558cbd1320e67aa8d488770583e 
 Final POST parameters: 
 - api_key: 1234 
- timestamp: 1315060510 
- public_id: sample_image 
- eager: w_400,h_300,c_pad|w_260,h_200,c_crop 
- signature: bfd09f95f331f558cbd1320e67aa8d488770583e 
- file: [your file] 
 Remember: exclude file, cloud_name, resource_type, and api_key from signature generation. (https://cloudinary.com/documentation/authentication_signatures) 
 
Source:
Generating authentication signatures | Documentation (https://cloudinary.com/documentation/authentication_signatures)

11:01 AM | Ask Cloudinary AI: Is that what you were looking for?

11:03 AM | Olivier Coq: yes. define: eager

11:03 AM | Ask Cloudinary AI: The eager parameter creates transformed versions of your uploaded asset immediately during upload, rather than generating them on-demand when first requested. (https://cloudinary.com/documentation/image_upload_api_reference) 
 What it does: 
 - Generates specific transformations right after upload 
- Accepts single transformations or pipe-separated lists (for REST API) 
- Useful for "cache warm up" to avoid delays on first access (https://cloudinary.com/documentation/eager_and_incoming_transformations) 
 Example formats: 
 - Single: w_400,h_300,c_pad 
- Multiple: w_400,h_300,c_pad|w_260,h_200,c_crop (https://cloudinary.com/documentation/upload_parameters) 
 This pre-generates the transformed images so they're ready immediately when your users need them, instead of processing transformations on-the-fly during delivery. 
 
Sources:
Upload API Reference | Documentation (https://cloudinary.com/documentation/image_upload_api_reference)
Eager and incoming transformations | Documentation (https://cloudinary.com/documentation/eager_and_incoming_transformations)
Customizing uploads | Documentation (https://cloudinary.com/documentation/upload_parameters)

11:03 AM | Ask Cloudinary AI: Did that answer your question?

11:07 AM | Ask Cloudinary AI: If you still need help with your signed upload or any other Cloudinary feature, I'm here to assist! Would you like to share more details about what you're trying to resolve or understand?

11:11 AM | Olivier Coq: one second

11:11 AM | Ask Cloudinary AI: Take your time! I'm here whenever you're ready to continue. 😊

11:16 AM | Ask Cloudinary AI: If you still need help with your signed upload or any other Cloudinary feature, I'm here to assist! Would you like to share more details about what you're trying to resolve or understand?

---
Exported from Cloudinary on October 30, 2025 at 01:40 PM Central Time (US & Canada) time CDT (GMT-0500)