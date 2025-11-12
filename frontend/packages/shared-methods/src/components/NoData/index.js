import React from 'react';
import Image from 'next/image';
import NodataSvg from '@/assets/images/nodata.svg';
export const NoData = ({ url, text = 'No Data', width, height, className }) => {
	return (
		<div
			className={`h-full w-full flex flex-col justify-center items-center py-5 ${className}`}
			style={{ width, height }}
		>
			{url ? (
				<Image src={url} alt="" width={120} height={120} />
			) : (
				<NodataSvg className="w-[120px] h-[120px]" />
			)}
			<span className="text-base-label">{text}</span>
		</div>
	);
};
