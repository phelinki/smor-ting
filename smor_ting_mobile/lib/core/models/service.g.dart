// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'service.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ServiceCategory _$ServiceCategoryFromJson(Map<String, dynamic> json) =>
    ServiceCategory(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String,
      icon: json['icon'] as String,
      color: json['color'] as String,
      isActive: json['is_active'] as bool,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );

Map<String, dynamic> _$ServiceCategoryToJson(ServiceCategory instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'description': instance.description,
      'icon': instance.icon,
      'color': instance.color,
      'is_active': instance.isActive,
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };

Service _$ServiceFromJson(Map<String, dynamic> json) => Service(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String,
      categoryId: json['category_id'] as String,
      providerId: json['provider_id'] as String,
      price: (json['price'] as num).toDouble(),
      currency: json['currency'] as String,
      duration: (json['duration'] as num).toInt(),
      images:
          (json['images'] as List<dynamic>).map((e) => e as String).toList(),
      isActive: json['is_active'] as bool,
      rating: (json['rating'] as num).toDouble(),
      reviewCount: (json['review_count'] as num).toInt(),
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );

Map<String, dynamic> _$ServiceToJson(Service instance) => <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'description': instance.description,
      'category_id': instance.categoryId,
      'provider_id': instance.providerId,
      'price': instance.price,
      'currency': instance.currency,
      'duration': instance.duration,
      'images': instance.images,
      'is_active': instance.isActive,
      'rating': instance.rating,
      'review_count': instance.reviewCount,
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };

ServiceProvider _$ServiceProviderFromJson(Map<String, dynamic> json) =>
    ServiceProvider(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      businessName: json['business_name'] as String,
      description: json['description'] as String,
      experience: (json['experience'] as num).toInt(),
      certifications: (json['certifications'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      serviceAreas: (json['service_areas'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      isVerified: json['is_verified'] as bool,
      rating: (json['rating'] as num).toDouble(),
      reviewCount: (json['review_count'] as num).toInt(),
      completedJobs: (json['completed_jobs'] as num).toInt(),
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );

Map<String, dynamic> _$ServiceProviderToJson(ServiceProvider instance) =>
    <String, dynamic>{
      'id': instance.id,
      'user_id': instance.userId,
      'business_name': instance.businessName,
      'description': instance.description,
      'experience': instance.experience,
      'certifications': instance.certifications,
      'service_areas': instance.serviceAreas,
      'is_verified': instance.isVerified,
      'rating': instance.rating,
      'review_count': instance.reviewCount,
      'completed_jobs': instance.completedJobs,
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };

Booking _$BookingFromJson(Map<String, dynamic> json) => Booking(
      id: json['id'] as String,
      customerId: json['customer_id'] as String,
      providerId: json['provider_id'] as String,
      serviceId: json['service_id'] as String,
      status: $enumDecode(_$BookingStatusEnumMap, json['status']),
      scheduledDate: DateTime.parse(json['scheduled_date'] as String),
      completedDate: json['completed_date'] == null
          ? null
          : DateTime.parse(json['completed_date'] as String),
      address: Address.fromJson(json['address'] as Map<String, dynamic>),
      notes: json['notes'] as String,
      totalAmount: (json['total_amount'] as num).toDouble(),
      currency: json['currency'] as String,
      paymentStatus: json['payment_status'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );

Map<String, dynamic> _$BookingToJson(Booking instance) => <String, dynamic>{
      'id': instance.id,
      'customer_id': instance.customerId,
      'provider_id': instance.providerId,
      'service_id': instance.serviceId,
      'status': _$BookingStatusEnumMap[instance.status]!,
      'scheduled_date': instance.scheduledDate.toIso8601String(),
      'completed_date': instance.completedDate?.toIso8601String(),
      'address': instance.address,
      'notes': instance.notes,
      'total_amount': instance.totalAmount,
      'currency': instance.currency,
      'payment_status': instance.paymentStatus,
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };

const _$BookingStatusEnumMap = {
  BookingStatus.pending: 'pending',
  BookingStatus.confirmed: 'confirmed',
  BookingStatus.inProgress: 'in_progress',
  BookingStatus.completed: 'completed',
  BookingStatus.cancelled: 'cancelled',
};

Address _$AddressFromJson(Map<String, dynamic> json) => Address(
      street: json['street'] as String,
      city: json['city'] as String,
      county: json['county'] as String,
      country: json['country'] as String,
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
    );

Map<String, dynamic> _$AddressToJson(Address instance) => <String, dynamic>{
      'street': instance.street,
      'city': instance.city,
      'county': instance.county,
      'country': instance.country,
      'latitude': instance.latitude,
      'longitude': instance.longitude,
    };

Review _$ReviewFromJson(Map<String, dynamic> json) => Review(
      id: json['id'] as String,
      bookingId: json['booking_id'] as String,
      customerId: json['customer_id'] as String,
      providerId: json['provider_id'] as String,
      serviceId: json['service_id'] as String,
      rating: (json['rating'] as num).toInt(),
      comment: json['comment'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$ReviewToJson(Review instance) => <String, dynamic>{
      'id': instance.id,
      'booking_id': instance.bookingId,
      'customer_id': instance.customerId,
      'provider_id': instance.providerId,
      'service_id': instance.serviceId,
      'rating': instance.rating,
      'comment': instance.comment,
      'created_at': instance.createdAt.toIso8601String(),
    };
