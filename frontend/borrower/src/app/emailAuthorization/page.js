'use client';
import React, { useState } from 'react';

// 授权gmail，以用户邮箱发送邮件<中转页,接收参数>
const EmailAuthorization = () => {
	const [isResult, setIsResult] = useState(false);
	const [result, setResult] = useState({});

	return <div>授权</div>;
};
export default EmailAuthorization;
