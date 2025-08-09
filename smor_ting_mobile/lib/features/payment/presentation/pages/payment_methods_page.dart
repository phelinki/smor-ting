import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

class PaymentMethod {
  final String id;
  final String name;
  final String lastFourDigits;
  final String type;
  final String icon;
  bool isDefault;

  PaymentMethod({
    required this.id,
    required this.name,
    required this.lastFourDigits,
    required this.type,
    required this.icon,
    this.isDefault = false,
  });
}

class PaymentMethodsPage extends ConsumerStatefulWidget {
  const PaymentMethodsPage({super.key});

  @override
  ConsumerState<PaymentMethodsPage> createState() => _PaymentMethodsPageState();
}

class _PaymentMethodsPageState extends ConsumerState<PaymentMethodsPage> {
  List<PaymentMethod> paymentMethods = [
    PaymentMethod(
      id: '1',
      name: 'Visa Card',
      lastFourDigits: '1234',
      type: 'Credit Card',
      icon: '💳',
      isDefault: true,
    ),
    PaymentMethod(
      id: '2',
      name: 'Mobile Money',
      lastFourDigits: 'MTN',
      type: 'Mobile Money',
      icon: '📱',
    ),
    PaymentMethod(
      id: '3',
      name: 'Mastercard',
      lastFourDigits: '5678',
      type: 'Credit Card',
      icon: '💳',
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF2F2F7),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF007AFF)),
          onPressed: () => context.pop(),
        ),
        title: const Text(
          'Payment Methods',
          style: TextStyle(
            color: Colors.black,
            fontSize: 18,
            fontWeight: FontWeight.w600,
          ),
        ),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.add, color: Color(0xFF007AFF)),
            onPressed: () {
              _showAddPaymentMethodDialog();
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Header Section
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(16),
            color: Colors.white,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Saved Payment Methods',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: Colors.black,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  'Manage your payment methods for quick and secure transactions',
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.grey[600],
                  ),
                ),
              ],
            ),
          ),
          
          // Payment Methods List
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: paymentMethods.length,
              itemBuilder: (context, index) {
                final method = paymentMethods[index];
                return _buildPaymentMethodCard(method);
              },
            ),
          ),
          
          // Add New Method Button
          Container(
            padding: const EdgeInsets.all(16),
            child: SizedBox(
              width: double.infinity,
              height: 56,
              child: ElevatedButton.icon(
                onPressed: () {
                  _showAddPaymentMethodDialog();
                },
                icon: const Icon(Icons.add, size: 20),
                label: const Text('Add New Payment Method'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF007AFF),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  elevation: 0,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPaymentMethodCard(PaymentMethod method) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: method.isDefault ? const Color(0xFFFF3B30) : const Color(0xFFE5E5EA),
          width: method.isDefault ? 2 : 1,
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: ListTile(
        contentPadding: const EdgeInsets.all(16),
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: const Color(0xFFF2F2F7),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Center(
            child: Text(
              method.icon,
              style: const TextStyle(fontSize: 24),
            ),
          ),
        ),
        title: Row(
          children: [
            Text(
              method.name,
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.black,
              ),
            ),
            if (method.isDefault) ...[
              const SizedBox(width: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: const Color(0xFFFF3B30),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Text(
                  'Default',
                  style: TextStyle(
                    fontSize: 10,
                    color: Colors.white,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
            ],
          ],
        ),
        subtitle: Text(
          '${method.type} •••• ${method.lastFourDigits}',
          style: const TextStyle(
            fontSize: 14,
            color: Color(0xFF8E8E93),
          ),
        ),
        trailing: PopupMenuButton<String>(
          icon: const Icon(Icons.more_vert, color: Color(0xFF8E8E93)),
          onSelected: (value) {
            _handlePaymentMethodAction(value, method);
          },
          itemBuilder: (context) => [
            if (!method.isDefault)
              const PopupMenuItem(
                value: 'set_default',
                child: Row(
                  children: [
                    Icon(Icons.star, color: Color(0xFF007AFF)),
                    SizedBox(width: 8),
                    Text('Set as Default'),
                  ],
                ),
              ),
            const PopupMenuItem(
              value: 'edit',
              child: Row(
                children: [
                  Icon(Icons.edit, color: Color(0xFF007AFF)),
                  SizedBox(width: 8),
                  Text('Edit'),
                ],
              ),
            ),
            const PopupMenuItem(
              value: 'delete',
              child: Row(
                children: [
                  Icon(Icons.delete, color: Color(0xFFFF3B30)),
                  SizedBox(width: 8),
                  Text('Delete'),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _handlePaymentMethodAction(String action, PaymentMethod method) {
    switch (action) {
      case 'set_default':
        setState(() {
          for (var pm in paymentMethods) {
            pm.isDefault = false;
          }
          method.isDefault = true;
        });
        _showSnackBar('${method.name} set as default');
        break;
      case 'edit':
        _showEditPaymentMethodDialog(method);
        break;
      case 'delete':
        _showDeleteConfirmationDialog(method);
        break;
    }
  }

  void _showAddPaymentMethodDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Add Payment Method'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              _buildPaymentOption('Credit/Debit Card', '💳', () {
                Navigator.of(context).pop();
                _showCardForm();
              }),
              const SizedBox(height: 12),
              _buildPaymentOption('Mobile Money', '📱', () {
                Navigator.of(context).pop();
                _showMobileMoneyForm();
              }),
              const SizedBox(height: 12),
              _buildPaymentOption('Bank Transfer', '🏦', () {
                Navigator.of(context).pop();
                _showBankForm();
              }),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
          ],
        );
      },
    );
  }

  Widget _buildPaymentOption(String title, String icon, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          border: Border.all(color: const Color(0xFFE5E5EA)),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            Text(icon, style: const TextStyle(fontSize: 24)),
            const SizedBox(width: 12),
            Text(
              title,
              style: const TextStyle(
                fontSize: 16,
                color: Colors.black,
              ),
            ),
            const Spacer(),
            const Icon(Icons.arrow_forward_ios, color: Color(0xFF8E8E93), size: 16),
          ],
        ),
      ),
    );
  }

  void _showCardForm() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Add Card'),
          content: const Text('Card form will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Card added successfully');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Add'),
            ),
          ],
        );
      },
    );
  }

  void _showMobileMoneyForm() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Add Mobile Money'),
          content: const Text('Mobile money form will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Mobile money added successfully');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Add'),
            ),
          ],
        );
      },
    );
  }

  void _showBankForm() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Add Bank Account'),
          content: const Text('Bank form will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Bank account added successfully');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Add'),
            ),
          ],
        );
      },
    );
  }

  void _showEditPaymentMethodDialog(PaymentMethod method) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text('Edit ${method.name}'),
          content: const Text('Edit form will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Payment method updated');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Update'),
            ),
          ],
        );
      },
    );
  }

  void _showDeleteConfirmationDialog(PaymentMethod method) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Delete Payment Method'),
          content: Text('Are you sure you want to delete ${method.name}?'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                setState(() {
                  paymentMethods.removeWhere((pm) => pm.id == method.id);
                });
                Navigator.of(context).pop();
                _showSnackBar('Payment method deleted');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Delete'),
            ),
          ],
        );
      },
    );
  }

  void _showSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: const Color(0xFF34C759),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
    );
  }
} 