import React from 'react';
import Image from 'next/image';
import NodataSvg from '@/assets/images/nodata.svg';
export const NoData = ({ url, text = 'No Data', width, height, className, children }) => {
	return (
		<div
			className={`h-full w-full flex flex-col justify-center items-center pb-5 pt-5 ${className}`}
			style={{ width, height }}
		>
			{children ? (
				children
			) : (
				<>
					{url ? (
						<Image src={url} alt="" width={120} height={120} />
					) : (
						<NodataSvg style={{ width: 120, height: 120 }} />
					)}
					<span className="text-base-label">{text}</span>
				</>
			)}
		</div>
	);
};
