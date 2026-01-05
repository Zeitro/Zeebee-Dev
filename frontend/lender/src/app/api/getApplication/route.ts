import { NextResponse } from 'next/server';
import path from 'path';
import fs from 'fs/promises';

export async function GET() {
	try {
		const pdfPath = path.join(process.cwd(), 'public', '/Application_Form.pdf');
		console.log(pdfPath);

		const pdfBuffer = await fs.readFile(pdfPath);

		// ✅ 核心改动：显式地从 Node.js Buffer 创建一个标准的 Uint8Array
		// 这一步创建了一个新的、类型明确的数组视图，消除了所有歧义。
		const uint8Array = new Uint8Array(pdfBuffer);

		// 现在 new Blob 接收的是一个它完全理解的类型
		const blob = new Blob([uint8Array]);

		return new NextResponse(blob, {
			status: 200,
			headers: {
				'Content-Type': 'application/pdf',
			},
		});
	} catch (error) {
		console.error('Failed to read PDF file:', error);
		return new NextResponse('Error finding PDF', { status: 404 });
	}
}
