import React from 'react';
import Image from 'next/image';
const NoData = ({ url, text = 'No Data', width, height }) => {
	return (
		<div
			className="h-full w-full flex flex-col justify-center items-center"
			style={{ width, height }}
		>
			<Image
				src={url || require('@/assets/images/nodata.svg')}
				alt=""
				className="w-[120px] h-[120px]"
			/>
			<span className='text-base-label'>{text}</span>
		</div>
	);
};
export default NoData;
