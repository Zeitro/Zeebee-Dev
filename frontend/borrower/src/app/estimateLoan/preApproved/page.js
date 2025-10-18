'use client';
import { Button, Form, Steps } from 'antd';
import React, { useEffect, useMemo, useState } from 'react';
// import styles from './index.module.scss';
import { useRouter } from 'next/navigation';
import { createRequest } from '@/lib/fetch';

// 预审计算loan
const PreApproved = () => {
	return (
		<div
			className={[
				' max-[676px] mt-10 bg-white rounded-lg px-8 py-6 h-max max-h-full overflow-hidden overflow-y-auto',
			].join(' ')}
		>
			<div className="text-3xl/10 font-medium text-[#0958D9]">Pre - Approved!</div>
		</div>
	);
};
export default PreApproved;
