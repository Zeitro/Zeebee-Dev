'use client';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import { AutoComplete, Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';
/**
 * Google Api
 * 根据输入内容，获取地址信息，并拿到city/county/state/zip/unit相关数据
 * 支持Form中使用
 * @param {*string} country us美国
 * @param {*function} onPlaceSelected 选中下拉某项后触发，返回city/county/state/zip/unit
 */
const GoogleAutoComplete = ({
	country = 'us',
	value,
	onChange = () => {},
	onPlaceSelected,
	placeholder = 'Enter address',
	...props
}) => {
	const [address, setAddress] = useState('');
	const [options, setOptions] = useState([]);
	const [loading, setLoading] = useState(false);
	const autocompleteService = useRef(null);
	const sessionToken = useRef(null);

	useEffect(() => {
		setAddress(value);
	}, [value]);

	const initAutocompleteService = () => {
		if (window.google?.maps?.places) {
			autocompleteService.current = new window.google.maps.places.AutocompleteService();
			sessionToken.current = new window.google.maps.places.AutocompleteSessionToken();
		}
	};

	const handleSearch = useCallback(
		(value) => {
			if (!autocompleteService.current || !value) {
				setOptions([]);
				return;
			}

			setLoading(true);
			autocompleteService.current.getPlacePredictions(
				{
					input: value,
					componentRestrictions: { country },
					sessionToken: sessionToken.current,
					types: ['address'],
				},
				(predictions, status) => {
					setLoading(false);
					if (status === 'OK' && predictions) {
						setOptions(
							predictions.map((p) => ({
								value: p.description,
								label: p.description,
								placeId: p.place_id,
							})),
						);
					} else {
						setOptions([]);
					}
				},
			);
		},
		[autocompleteService, country, setOptions, setLoading, sessionToken],
	);

	const handleSelect = async (value, option) => {
		if (!window.google.maps.places) return;
		const placesService = new window.google.maps.places.PlacesService(
			document.createElement('div'),
		);

		placesService.getDetails(
			{
				placeId: option.placeId,
				fields: ['address_components', 'geometry', 'formatted_address'],
				sessionToken: sessionToken.current,
			},
			(place, status) => {
				if (status === 'OK') {
					const addressInfo = parseAddressComponents(place.address_components);
					const address = addressInfo?.streetNumber
						? `${addressInfo?.streetNumber} ${addressInfo?.street}`
						: addressInfo?.street;
					onChange(address);
					setAddress(address);
					const zip = addressInfo?.zipCode_suffix
						? `${addressInfo?.zipCode}-${addressInfo?.zipCode_suffix}`
						: addressInfo?.zipCode || '';
					onPlaceSelected?.({
						...addressInfo,
						location: place.geometry?.location,
						fullAddress: place.formatted_address,
						zip,
					});
				}
				sessionToken.current = new window.google.maps.places.AutocompleteSessionToken();
			},
		);
	};

	// 获取city/unit/state/county/zip
	const parseAddressComponents = (components) => {
		const result = {};
		components.forEach((component) => {
			if (component.types.includes('locality')) {
				result.city = component.long_name;
			} else if (component.types.includes('administrative_area_level_1')) {
				result.state = component.short_name;
			} else if (component.types.includes('postal_code')) {
				result.zipCode = component.long_name;
			} else if (component.types.includes('administrative_area_level_2')) {
				result.country = component.long_name;
			} else if (component.types.includes('route')) {
				result.street = component.short_name;
			} else if (component.types.includes('street_number')) {
				result.streetNumber = component.long_name;
			} else if (component.types.includes('subpremise')) {
				result.unit = component.long_name;
			} else if (component.types.includes('postal_code_suffix')) {
				result.zipCode_suffix = component.short_name;
			}
		});
		return result;
	};

	useEffect(() => {
		if (window.google) {
			initAutocompleteService();
		}
	}, []);
	return (
		<AutoComplete
			value={address}
			options={options}
			onChange={(e) => {
				setAddress(e);
				onChange(e);
			}}
			onSearch={handleSearch}
			onSelect={handleSelect}
			placeholder={placeholder}
			notFoundContent={loading ? <Spin indicator={<LoadingOutlined spin />} /> : null}
			style={{ width: '100%' }}
			{...props}
		/>
	);
};

export default GoogleAutoComplete;
