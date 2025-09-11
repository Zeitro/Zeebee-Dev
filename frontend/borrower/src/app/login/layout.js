import Image from 'next/image';
import React from 'react';
const LoginLayout = ({ children }) => {
	return (
		<div className="h-full w-full flex flex-col">
			<div className="h-20 flex-shrink-0">
				<Image
					src="/images/logo.png"
					width={145}
					height={80}
					alt="logo"
					className="!h-full"
				/>
			</div>
			{children}
		</div>
	);
};
export default LoginLayout;
