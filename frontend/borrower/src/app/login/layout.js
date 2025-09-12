import Image from 'next/image';
import React from 'react';
const LoginLayout = ({ children }) => {
	return (
		<div className="h-full w-full flex flex-col">
			<div className="h-20 flex-shrink-0 flex items-center pl-6">
				<div className="h-7">
					<Image
						src="/images/logo.svg"
						width={95}
						height={30}
						alt="logo"
						className="!h-full"
					/>
				</div>
			</div>
			{children}
		</div>
	);
};
export default LoginLayout;
