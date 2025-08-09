import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/connectivity_provider.dart';


import '../../../../core/theme/app_theme.dart';
import '../../../../core/services/api_service.dart';
import '../../../../core/services/wallet_cache.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class WalletPage extends ConsumerStatefulWidget {
  const WalletPage({super.key});

  @override
  ConsumerState<WalletPage> createState() => _WalletPageState();
}

class _WalletPageState extends ConsumerState<WalletPage> {
  int _selectedTabIndex = 0;
  late final ApiService _api;
  late final WalletCache _cache;
  Map<String, dynamic>? _balances;

  @override
  void initState() {
    super.initState();
    _api = ref.read(apiServiceProvider);
    _cache = WalletCache(const FlutterSecureStorage());
    _cache.init().then((_) {
      setState(() {
        _balances = _cache.getBalances();
      });
      _refreshBalancesIfOnline();
    });
  }

  Future<void> _refreshBalancesIfOnline() async {
    final net = ref.read(connectivityProvider);
    if (!net.isOnline) return;
    try {
      final b = await _api.getWalletBalances();
      await _cache.saveBalances(b);
      if (mounted) {
        setState(() {
          _balances = b;
        });
      }
    } catch (_) {}
  }

  @override
  Widget build(BuildContext context) {
    final net = ref.watch(connectivityProvider);
    return Scaffold(
      backgroundColor: AppTheme.white,
      appBar: AppBar(
        backgroundColor: AppTheme.white,
        elevation: 0,
        title: Text(
          'Wallet',
          style: Theme.of(context).textTheme.titleLarge?.copyWith(
            fontWeight: FontWeight.bold,
            color: AppTheme.textPrimary,
          ),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.history, color: AppTheme.secondaryBlue),
            onPressed: () {
              // TODO: Navigate to full transaction history
            },
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Balance Card
            Container(
              margin: const EdgeInsets.all(16),
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                gradient: const LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    AppTheme.secondaryBlue,
                    Color(0xFF0056CC),
                  ],
                ),
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(
                    color: AppTheme.secondaryBlue.withValues(alpha: 0.3),
                    blurRadius: 15,
                    offset: const Offset(0, 8),
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Available Balance',
                        style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          color: AppTheme.white.withValues(alpha: 0.9),
                        ),
                      ),
                      const Icon(
                        Icons.account_balance_wallet,
                        color: AppTheme.white,
                        size: 24,
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                   Text(
                     net.isOnline
                         ? (_balances != null ? '${_balances!['available']} ${_balances!['currency']}' : '...')
                         : 'Offline',
                    style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                      color: AppTheme.white,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: AppTheme.white.withValues(alpha: 0.2),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(
                              Icons.trending_up,
                              color: AppTheme.successGreen,
                              size: 16,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              '+\$25.00 this week',
                              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: AppTheme.white,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            
            // Action Buttons
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  Expanded(
                    child: _buildActionButton(
                      icon: Icons.add,
                      label: 'Add Money',
                      color: AppTheme.primaryRed,
                      onTap: net.isOnline
                          ? () => _showAddMoneyDialog()
                          : () => _showOfflineSnack(),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildActionButton(
                      icon: Icons.send,
                      label: 'Send Money',
                      color: AppTheme.secondaryBlue,
                      onTap: net.isOnline
                          ? () => _showSendMoneyDialog()
                          : () => _showOfflineSnack(),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildActionButton(
                      icon: Icons.receipt_long,
                      label: 'Transactions',
                      color: AppTheme.successGreen,
                      onTap: () => _showTransactionHistory(),
                    ),
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 32),
            
            // Tab Bar
            Container(
              margin: const EdgeInsets.symmetric(horizontal: 16),
              decoration: BoxDecoration(
                color: AppTheme.lightGray,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Row(
                children: [
                  Expanded(
                    child: _buildTabButton(
                      title: 'Payment Methods',
                      index: 0,
                      isSelected: _selectedTabIndex == 0,
                    ),
                  ),
                  Expanded(
                    child: _buildTabButton(
                      title: 'Recent Activity',
                      index: 1,
                      isSelected: _selectedTabIndex == 1,
                    ),
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 16),
            
            // Tab Content
            if (_selectedTabIndex == 0) _buildPaymentMethodsTab(),
            if (_selectedTabIndex == 1) _buildRecentActivityTab(),
            
            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }

  Widget _buildActionButton({
    required IconData icon,
    required String label,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        decoration: BoxDecoration(
          color: color.withValues(alpha: 0.1),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(
              icon,
              color: color,
              size: 24,
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                color: color,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTabButton({
    required String title,
    required int index,
    required bool isSelected,
  }) {
    return GestureDetector(
      onTap: () {
        setState(() {
          _selectedTabIndex = index;
        });
      },
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: isSelected ? AppTheme.white : Colors.transparent,
          borderRadius: BorderRadius.circular(10),
        ),
        child: Text(
          title,
          textAlign: TextAlign.center,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            color: isSelected ? AppTheme.textPrimary : AppTheme.textSecondary,
            fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
          ),
        ),
      ),
    );
  }

  Widget _buildPaymentMethodsTab() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Saved Payment Methods
          ..._samplePaymentMethods.map((method) => _buildPaymentMethodCard(method)),
          
          const SizedBox(height: 16),
          
          // Add New Payment Method
          GestureDetector(
            onTap: () => _showAddPaymentMethodDialog(),
            child: Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                border: Border.all(
                  color: AppTheme.secondaryBlue,
                  width: 2,
                ),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Row(
                children: [
                  const Icon(
                    Icons.add_circle_outline,
                    color: AppTheme.secondaryBlue,
                    size: 24,
                  ),
                  const SizedBox(width: 12),
                  Text(
                    'Add New Payment Method',
                    style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                      color: AppTheme.secondaryBlue,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
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
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          Container(
            width: 48,
            height: 32,
            decoration: BoxDecoration(
              color: AppTheme.lightGray,
              borderRadius: BorderRadius.circular(8),
            ),
            child: Icon(
              method.type == 'card' ? Icons.credit_card : Icons.account_balance,
              color: AppTheme.secondaryBlue,
              size: 20,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  method.name,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                Text(
                  method.details,
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: AppTheme.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          if (method.isDefault)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: AppTheme.successGreen.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(
                'Default',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: AppTheme.successGreen,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildRecentActivityTab() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        children: _sampleTransactions.map((transaction) => _buildTransactionItem(transaction)).toList(),
      ),
    );
  }

  Widget _buildTransactionItem(Transaction transaction) {
    final isCredit = transaction.type == 'credit';
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppTheme.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: (isCredit ? AppTheme.successGreen : AppTheme.primaryRed).withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Icon(
              isCredit ? Icons.arrow_downward : Icons.arrow_upward,
              color: isCredit ? AppTheme.successGreen : AppTheme.primaryRed,
              size: 20,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  transaction.description,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppTheme.textPrimary,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  transaction.date,
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: AppTheme.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          Text(
            '${isCredit ? '+' : '-'}\$${transaction.amount.toStringAsFixed(2)}',
            style: Theme.of(context).textTheme.titleSmall?.copyWith(
              color: isCredit ? AppTheme.successGreen : AppTheme.primaryRed,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  void _showAddMoneyDialog() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _buildAddMoneySheet(),
    );
  }

  void _showSendMoneyDialog() {
    // TODO: Implement send money dialog
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Send money feature coming soon!'),
        backgroundColor: AppTheme.secondaryBlue,
      ),
    );
  }

  void _showTransactionHistory() {
    // TODO: Navigate to full transaction history
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Full transaction history coming soon!'),
        backgroundColor: AppTheme.secondaryBlue,
      ),
    );
  }

  void _showAddPaymentMethodDialog() {
    // TODO: Implement add payment method dialog
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Add payment method feature coming soon!'),
        backgroundColor: AppTheme.secondaryBlue,
      ),
    );
  }

  void _showOfflineSnack() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Wallet actions require internet connection.'),
        backgroundColor: AppTheme.primaryRed,
      ),
    );
  }

  Widget _buildAddMoneySheet() {
    return Container(
      padding: EdgeInsets.only(
        top: 24,
        left: 24,
        right: 24,
        bottom: MediaQuery.of(context).viewInsets.bottom + 24,
      ),
      decoration: const BoxDecoration(
        color: AppTheme.white,
        borderRadius: BorderRadius.only(
          topLeft: Radius.circular(24),
          topRight: Radius.circular(24),
        ),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Handle bar
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: AppTheme.textSecondary.withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          
          const SizedBox(height: 24),
          
          Text(
            'Add Money',
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: AppTheme.textPrimary,
            ),
            textAlign: TextAlign.center,
          ),
          
          const SizedBox(height: 24),
          
          // Quick Amount Buttons
          Row(
            children: [
              Expanded(child: _buildQuickAmountButton(25)),
              const SizedBox(width: 8),
              Expanded(child: _buildQuickAmountButton(50)),
              const SizedBox(width: 8),
              Expanded(child: _buildQuickAmountButton(100)),
              const SizedBox(width: 8),
              Expanded(child: _buildQuickAmountButton(200)),
            ],
          ),
          
          const SizedBox(height: 24),
          
          // Custom Amount Input
          TextField(
            keyboardType: TextInputType.number,
            decoration: InputDecoration(
              labelText: 'Custom Amount',
              hintText: 'Enter amount',
              prefixText: '\$',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
          
          const SizedBox(height: 24),
          
          // Add Money Button
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Money added successfully!'),
                  backgroundColor: AppTheme.successGreen,
                ),
              );
            },
            child: const Text('Add Money'),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickAmountButton(int amount) {
    return GestureDetector(
      onTap: () {
        // TODO: Set amount in text field
      },
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          border: Border.all(color: AppTheme.secondaryBlue),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(
          '\$$amount',
          textAlign: TextAlign.center,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            color: AppTheme.secondaryBlue,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
    );
  }
}

// Data classes
class PaymentMethod {
  final String id;
  final String name;
  final String details;
  final String type;
  final bool isDefault;

  PaymentMethod({
    required this.id,
    required this.name,
    required this.details,
    required this.type,
    required this.isDefault,
  });
}

class Transaction {
  final String id;
  final String description;
  final String date;
  final double amount;
  final String type; // 'credit' or 'debit'

  Transaction({
    required this.id,
    required this.description,
    required this.date,
    required this.amount,
    required this.type,
  });
}

// Sample data
final List<PaymentMethod> _samplePaymentMethods = [
  PaymentMethod(
    id: '1',
    name: 'Visa •••• 4242',
    details: 'Expires 12/25',
    type: 'card',
    isDefault: true,
  ),
  PaymentMethod(
    id: '2',
    name: 'Bank of America',
    details: 'Checking •••• 1234',
    type: 'bank',
    isDefault: false,
  ),
  PaymentMethod(
    id: '3',
    name: 'PayPal',
    details: 'john@example.com',
    type: 'paypal',
    isDefault: false,
  ),
];

final List<Transaction> _sampleTransactions = [
  Transaction(
    id: '1',
    description: 'Plumbing Service Payment',
    date: 'Today, 2:30 PM',
    amount: 85.00,
    type: 'debit',
  ),
  Transaction(
    id: '2',
    description: 'Wallet Top-up',
    date: 'Yesterday, 10:15 AM',
    amount: 100.00,
    type: 'credit',
  ),
  Transaction(
    id: '3',
    description: 'Cleaning Service Payment',
    date: 'Dec 15, 4:45 PM',
    amount: 60.00,
    type: 'debit',
  ),
  Transaction(
    id: '4',
    description: 'Refund - Cancelled Service',
    date: 'Dec 14, 11:20 AM',
    amount: 45.00,
    type: 'credit',
  ),
  Transaction(
    id: '5',
    description: 'Electrical Service Payment',
    date: 'Dec 12, 9:30 AM',
    amount: 120.00,
    type: 'debit',
  ),
];