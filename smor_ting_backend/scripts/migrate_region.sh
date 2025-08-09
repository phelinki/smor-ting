#!/bin/bash

# Smor-Ting MongoDB Atlas Region Migration Script
# This script helps migrate your cluster from one region to another

set -e

echo "🌍 MongoDB Atlas Region Migration"
echo "================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo ""
print_info "This script will help you migrate your MongoDB Atlas cluster to a different region."
echo ""

echo "📋 Migration Options:"
echo ""

echo "1. 🚀 Live Migration (Recommended)"
echo "   - Zero downtime"
echo "   - Automatic data sync"
echo "   - Requires M10+ cluster"
echo "   - Higher cost"
echo ""

echo "2. 📦 Export/Import Migration"
echo "   - Some downtime required"
echo "   - Manual process"
echo "   - Works with any cluster tier"
echo "   - Lower cost"
echo ""

echo "3. 🔄 Atlas Migration Tools"
echo "   - Semi-automated"
echo "   - Minimal downtime"
echo "   - Works with M2+ clusters"
echo "   - Medium cost"
echo ""

read -p "Choose migration method (1-3): " MIGRATION_METHOD

case $MIGRATION_METHOD in
    1)
        print_info "Live Migration (M10+ Required)"
        echo ""
        echo "Steps:"
        echo "1. Upgrade to M10+ cluster"
        echo "2. Go to Atlas Dashboard → Database → Your Cluster"
        echo "3. Click 'Configuration' tab"
        echo "4. Click 'Edit Configuration'"
        echo "5. Change region to South Africa"
        echo "6. Click 'Save Changes'"
        echo "7. Wait for migration to complete (2-4 hours)"
        echo ""
        print_warning "Note: This requires M10+ cluster ($57/month minimum)"
        ;;
    2)
        print_info "Export/Import Migration"
        echo ""
        echo "Steps:"
        echo "1. Create new cluster in South Africa"
        echo "2. Export data from current cluster"
        echo "3. Import data to new cluster"
        echo "4. Update connection string"
        echo "5. Test application"
        echo "6. Delete old cluster"
        echo ""
        print_warning "Note: Some downtime required during import"
        ;;
    3)
        print_info "Atlas Migration Tools"
        echo ""
        echo "Steps:"
        echo "1. Go to Atlas Dashboard → Data Migration"
        echo "2. Create new cluster in South Africa"
        echo "3. Use Atlas migration tools"
        echo "4. Monitor migration progress"
        echo "5. Update connection string"
        echo "6. Test application"
        echo "7. Delete old cluster"
        echo ""
        print_warning "Note: Requires M2+ cluster ($9/month minimum)"
        ;;
    *)
        print_error "Invalid choice. Please run the script again."
        exit 1
        ;;
esac

echo ""
print_info "📋 Pre-Migration Checklist:"
echo ""

echo "Before starting migration:"
echo "✅ Backup your current data"
echo "✅ Test your application thoroughly"
echo "✅ Plan maintenance window (if needed)"
echo "✅ Update your team about the migration"
echo "✅ Prepare rollback plan"
echo ""

print_info "🌍 Recommended Regions for Liberia:"
echo ""

echo "Primary Options:"
echo "1. South Africa (Johannesburg) - AWS"
echo "   - Latency to Liberia: ~200-300ms"
echo "   - Cost: Standard"
echo ""

echo "2. Europe (Ireland) - AWS"
echo "   - Latency to Liberia: ~150-200ms"
echo "   - Cost: Standard"
echo ""

echo "3. US East (N. Virginia) - AWS"
echo "   - Latency to Liberia: ~200-250ms"
echo "   - Cost: Standard"
echo ""

print_info "🔧 Post-Migration Steps:"
echo ""

echo "After migration:"
echo "1. Update connection string in .env file"
echo "2. Test all application features"
echo "3. Monitor performance metrics"
echo "4. Update DNS if using custom domain"
echo "5. Update mobile app configuration"
echo "6. Monitor user experience"
echo ""

print_info "📊 Performance Monitoring:"
echo ""

echo "Monitor these metrics after migration:"
echo "- Database response times"
echo "- Application latency"
echo "- User experience feedback"
echo "- Error rates"
echo "- Connection pool usage"
echo ""

print_status "Migration guide completed!"
echo ""
print_info "Next steps:"
echo "1. Choose your migration method"
echo "2. Follow the steps above"
echo "3. Test thoroughly after migration"
echo "4. Monitor performance"
echo ""

print_warning "⚠️  Important Notes:"
echo "- Always backup before migration"
echo "- Test in staging environment first"
echo "- Plan for potential downtime"
echo "- Have rollback plan ready"
echo "- Monitor closely after migration" 